package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
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

	r_p := new(Rep_Permission)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &r_p); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	if r_p.User_name == "" {
		return rsp.Json(400, ErrNoParameter("username"))
	}

	if r_p.Opt_permission != PERMISSION_READ && r_p.Opt_permission != PERMISSION_WRITE {
		return rsp.Json(400, ErrInvalidParameter("opt_permission"))
	}
	r_p.Repository_name = p.Repository_name

	Q := bson.M{COL_PERMIT_REPNAME: p.Repository_name, COL_PERMIT_USER: r_p.User_name}

	if _, err := db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Upsert(Q, r_p); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))

}

func getRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()

	Q := bson.M{COL_REPNAME: p.Repository_name}
	if user := strings.TrimSpace(r.FormValue("username")); user != "" {
		Q[COL_PERMIT_USER] = user
	}

	l, err := db.getPermits(C_REPOSITORY_PERMISSION, Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), l)
}

func getItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Item_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()

	Q := bson.M{COL_PERMIT_REPNAME: p.Repository_name, COL_PERMIT_ITEMNAME: p.Dataitem_name}
	if user := strings.TrimSpace(r.FormValue("username")); user != "" {
		Q[COL_PERMIT_USER] = user
	}

	l, err := db.getPermits(C_DATAITEM_PERMISSION, Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), l)
}

func delRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	defer db.Close()
	r.ParseForm()
	users := r.Form["username"]
	deleteAll := strings.TrimSpace(r.FormValue("delall"))
	if len(users) == 0 && deleteAll == "" {
		return rsp.Json(400, E(ErrorCodeNoParameter))
	}

	cmdCondiction := bson.M{CMD_IN: users}

	exec := bson.M{COL_REPNAME: p.Repository_name}
	if deleteAll != DELETE_PERMISSION_USR_ALL {
		exec[COL_PERMIT_USER] = cmdCondiction
	}

	if err := db.delPermit(C_REPOSITORY_PERMISSION, exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

func setItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Item_Permission) (int, string) {
	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err", err)
	}

	i_p := new(Item_Permission)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &i_p); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	if i_p.User_name == "" {
		return rsp.Json(400, ErrNoParameter("User_name"))
	}
	i_p.Dataitem_name = p.Dataitem_name
	i_p.Repository_name = p.Repository_name

	Q := bson.M{COL_PERMIT_REPNAME: p.Repository_name, COL_PERMIT_ITEMNAME: p.Dataitem_name, COL_PERMIT_USER: i_p.User_name}

	if _, err := db.DB(DB_NAME).C(C_DATAITEM_PERMISSION).Upsert(Q, i_p); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
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
