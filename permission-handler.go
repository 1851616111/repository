package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	PERMISSION_WRITE          = 1
	PERMISSION_READ           = 0
	DELETE_PERMISSION_USR_ALL = "1"
)

func setRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err", err)
	}

	result := new(Rep_Permission)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	if result.User_name == "" {
		return rsp.Json(400, ErrNoParameter("username"))
	}

	if result.Opt_permission != PERMISSION_READ && result.Opt_permission != PERMISSION_WRITE {
		return rsp.Json(400, ErrInvalidParameter("opt_permission"))
	}

	putRepositoryPermission(p.Repository_name, result.User_name, result.Opt_permission)

	//todo set to direct insert
	time.Sleep(time.Millisecond * 100)

	return rsp.Json(200, E(OK))
}

func getRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()
	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}

	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size < -1 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}

	Q := bson.M{COL_REPNAME: p.Repository_name}
	if user := strings.TrimSpace(r.FormValue("username")); user != "" {
		Q[COL_PERMIT_USER] = user
	}
	if cooperator := strings.TrimSpace(r.FormValue("cooperator")); cooperator != "" {
		Q["opt_permission"] = PERMISSION_WRITE
	}

	l, err := db.getPermits(C_REPOSITORY_PERMISSION, Q, []int{page_index, page_size})
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	if list, ok := l.([]Rep_Permission); ok {
		if len(list) == 0 {
			return rsp.Json(400, ErrQueryNotFound(""))
		}
	}

	n, _ := db.countPermits(C_REPOSITORY_PERMISSION, Q)
	res := struct {
		L     interface{} `json:"permissions"`
		Total int         `json:"total"`
	}{
		L:     l,
		Total: n,
	}
	return rsp.Json(200, E(OK), res)
}

func getItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Item_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()
	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}

	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size < -1 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}
	Q := bson.M{COL_PERMIT_REPNAME: p.Repository_name, COL_PERMIT_ITEMNAME: p.Dataitem_name}
	if user := strings.TrimSpace(r.FormValue("username")); user != "" {
		Q[COL_PERMIT_USER] = user
	}

	l, err := db.getPermits(C_DATAITEM_PERMISSION, Q, []int{page_index, page_size})
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	if list, ok := l.([]Item_Permission); ok {
		if len(list) == 0 {
			return rsp.Json(400, ErrQueryNotFound(""))
		}
	}

	n, _ := db.countPermits(C_DATAITEM_PERMISSION, Q)
	res := struct {
		L     interface{} `json:"permissions"`
		Total int         `json:"total"`
	}{
		L:     l,
		Total: n,
	}
	return rsp.Json(200, E(OK), res)
}

func delRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()
	deleteAll := strings.TrimSpace(r.FormValue("delall"))

	selector := bson.M{COL_REPNAME: p.Repository_name}
	if deleteAll != DELETE_PERMISSION_USR_ALL {
		selector[COL_PERMIT_USER] = p.User_name
	}

	Log.Infof(" %#v\n", selector)
	db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).RemoveAll(selector)

	return rsp.Json(200, E(OK))
}

func delRepCoptPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission, RepAccessType string) (int, string) {
	defer db.Close()
	r.ParseForm()
	deleteAll := strings.TrimSpace(r.FormValue("delall"))
	repName := p.Repository_name

	selector := bson.M{COL_REPNAME: repName}
	if deleteAll != DELETE_PERMISSION_USR_ALL {
		selector[COL_PERMIT_USER] = p.User_name
	}

	result := Rep_Permission{}
	execs := []Execute{}
	toDelUsers := []string{}

	iter := db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Find(selector).Iter()
	for iter.Next(&result) {

		switch RepAccessType {
		case ACCESS_PRIVATE:

			execs = []Execute{
				Execute{
					Collection: C_REPOSITORY_PERMISSION,
					Selector:   bson.M{COL_REPNAME: repName, COL_PERMIT_USER: result.User_name},
					Update:     bson.M{CMD_SET: bson.M{"opt_permission": 0}},
					Type:       Exec_Type_Update,
				},
				Execute{
					Collection: C_REPOSITORY,
					Selector:   bson.M{COL_REPNAME: repName},
					Update:     bson.M{CMD_PULL: bson.M{COL_REP_COOPERATOR: result.User_name}},
					Type:       Exec_Type_Update,
				},
			}

		case ACCESS_PUBLIC:
			toDelUsers = append(toDelUsers, result.User_name)
		}
	}

	if RepAccessType == ACCESS_PUBLIC && len(toDelUsers) > 0 {
		delete := bson.M{
			COL_REPNAME:     repName,
			COL_PERMIT_USER: bson.M{CMD_IN: toDelUsers},
		}
		db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Remove(delete)
	}

	go asynExec(execs...)

	return rsp.Json(200, E(OK))
}

func setItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, p Item_Permission) (int, string) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err", err)
	}

	result := new(Item_Permission)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	if result.User_name == "" {
		return rsp.Json(400, ErrNoParameter("User_name"))
	}
	putDataitemPermission(p.Repository_name, p.Dataitem_name, result.User_name)

	//todo set to direct insert
	time.Sleep(time.Millisecond * 100)

	return rsp.Json(200, E(OK))
}

func delItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Item_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()

	users := r.Form["username"]
	deleteAll := strings.TrimSpace(r.FormValue("delall"))

	if len(users) == 0 && deleteAll == "" {
		return rsp.Json(400, E(ErrorCodeNoParameter))
	}

	cmdCondiction := bson.M{CMD_IN: users}
	exec := bson.M{COL_REPNAME: p.Repository_name, COL_PERMIT_ITEMNAME: p.Dataitem_name}
	if deleteAll != DELETE_PERMISSION_USR_ALL {
		exec[COL_PERMIT_USER] = cmdCondiction
	}

	if err := db.delPermit(C_DATAITEM_PERMISSION, exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}
