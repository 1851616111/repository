package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

const (
	AUTHORIZATION = "Authorization"

	USER_SERVICE_RET_USERTYPE = "userType"
)

var (
	API_SERVER           = Env("API_SERVER", false)
	API_PORT             = Env("API_PORT", false)
	USER_TP_ADMIN    int = 2
	USER_TP_UNKNOW       = -1
	Unauthroized     string
	PermissionDenied string
)

func init() {
	b, _ := json.Marshal(E(ErrorCodeUnauthorized))
	Unauthroized = string(b)
	b, _ = json.Marshal(E(ErrorCodePermissionDenied))
	PermissionDenied = string(b)
}

func auth(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	if login_Name == "" {
		http.Error(w, Unauthroized, 401)
	}
	c.Map(login_Name)
	return

}

func getUserType(r *http.Request, db *DB) int {
	login_Name := r.Header.Get("User")
	token := r.Header.Get(AUTHORIZATION)
	if login_Name == "" || login_Name == "" {
		return USER_TP_UNKNOW
	}
	b, err := httpGet(fmt.Sprintf("http://%s:%s/users/%s", API_SERVER, API_PORT, login_Name), AUTHORIZATION, token)
	get(err)
	result := new(Result)
	err = json.Unmarshal(b, result)
	get(err)
	if result.Data != nil {
		u := result.Data.(map[string]interface{})
		if userType, exist := u[USER_SERVICE_RET_USERTYPE]; exist {
			return int(userType.(float64))
		}
	}
	return USER_TP_UNKNOW
}

func authAdmin(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {

	if getUserType(r, db) == USER_TP_ADMIN {
		login_Name := r.Header.Get("User")
		c.Map(login_Name)
		return
	} else {
		http.Error(w, Unauthroized, 401)
		return
	}
}

func chkRepPermission(w http.ResponseWriter, r *http.Request, param martini.Params, c martini.Context, db *DB) {
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, Unauthroized, 401)
		return
	}
	repName := strings.TrimSpace(param["repname"])
	if repName == "" {
		http.Error(w, ErrNoParameter("repname").ErrToString(), 401)
		return
	}

	if rep, _ := db.getRepository(bson.M{COL_REPNAME: repName}); rep.Create_user != user {
		http.Error(w, PermissionDenied, 401)
		return
	}
	c.Map(Rep_Permission{Repository_name: repName})
}

func chkItemPermission(w http.ResponseWriter, r *http.Request, param martini.Params, c martini.Context, db *DB) {
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, Unauthroized, 401)
		return
	}
	repName := strings.TrimSpace(param["repname"])
	if repName == "" {
		http.Error(w, ErrNoParameter("repname").ErrToString(), 401)
		return
	}
	itemname := strings.TrimSpace(param["itemname"])
	if repName == "" {
		http.Error(w, ErrNoParameter("itemname").ErrToString(), 401)
		return
	}

	if item, _ := db.getDataitem(bson.M{COL_REPNAME: repName, COL_ITEM_NAME: itemname}); item.Create_user != user {
		http.Error(w, PermissionDenied, 401)
		return
	}
	c.Map(Item_Permission{Repository_name: repName, Dataitem_name: itemname})
}
