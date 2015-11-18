package main

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

const (
	ADMIN = "admin"
	PANXY = "panxy3@asiainfo.com"
)

func auth(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	if login_Name == "" {
		http.Error(w, "unauthorized", 401)
	}
	c.Map(login_Name)
	return

}

func authAdmin(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	if login_Name != ADMIN && login_Name != PANXY {
		http.Error(w, "unauthorized", 401)
		return
	}
	c.Map(login_Name)
	return
}

func chkRepPermission(w http.ResponseWriter, r *http.Request, param martini.Params, c martini.Context, db *DB) {
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, "unauthorized", 401)
		return
	}
	repName := strings.TrimSpace(param["repname"])
	if repName == "" {
		http.Error(w, "no param repname", 401)
		return
	}

	if rep, _ := db.getRepository(bson.M{COL_REPNAME: repName}); rep.Create_user != user {
		http.Error(w, "no privilage", 401)
		return
	}
	c.Map(Rep_Permission{Repository_name: repName})
}

func chkItemPermission(w http.ResponseWriter, r *http.Request, param martini.Params, c martini.Context, db *DB) {
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, "unauthorized", 401)
		return
	}
	repName := strings.TrimSpace(param["repname"])
	if repName == "" {
		http.Error(w, "no param repname", 401)
		return
	}
	itemname := strings.TrimSpace(param["itemname"])
	if repName == "" {
		http.Error(w, "no param itemname", 401)
		return
	}

	if item, _ := db.getDataitem(bson.M{COL_REPNAME: repName, COL_ITEM_NAME: itemname}); item.Create_user != user {
		http.Error(w, "no privilage", 401)
		return
	}
	c.Map(Item_Permission{Dataitem_name: itemname})
}
