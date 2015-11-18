package main

import (
	"github.com/go-martini/martini"
	"net/http"
	//"io/ioutil"
	//"log"
	//"encoding/json"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"strings"
)

func setRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("read request body err", err)
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
	log.Println("-------------->", r_p.Write)
	if r_p.Write != 1 {
		return rsp.Json(400, ErrInvalidParameter("write"))
	}

	r_p.Repository_name = p.Repository_name

	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Insert(r_p); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))

}

func getRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	Q := bson.M{COL_REPNAME: p.Repository_name}
	l, err := db.getPermit(C_REPOSITORY_PERMISSION, Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), l)
}

func getItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	Q := bson.M{COL_REPNAME: p.Repository_name}
	l, err := db.getPermit(C_DATAITEM_PERMISSION, Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), l)
}

func delRepPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Rep_Permission) (int, string) {
	username := strings.TrimSpace(r.FormValue("username"))
	if username == "" {
		return rsp.Json(400, ErrNoParameter(username))
	}

	exec := bson.M{COL_REPNAME: p.Repository_name, COL_CREATE_USER: username}
	if err := db.delPermit(C_REPOSITORY_PERMISSION, exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

func setItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Item_Permission) (int, string) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("read request body err", err)
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

	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Insert(i_p); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}

func delItemPmsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, p Item_Permission) (int, string) {
	username := strings.TrimSpace(r.FormValue("username"))
	if username == "" {
		return rsp.Json(400, ErrNoParameter(username))
	}

	exec := bson.M{COL_REPNAME: p.Dataitem_name, COL_CREATE_USER: username}
	if err := db.delPermit(C_DATAITEM_PERMISSION, exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}
