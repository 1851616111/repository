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
	Quota_Rep_Pub             = "quotaPublic"
	Quota_Rep_Pri             = "quotaPrivate"
	VIP_SERVICE_ADMIN_PUB     = -1
	VIP_SERVICE_ADMIN_PRI     = -1
	Subscripters_By_Rep       = "repository"
	Subscripters_By_Item      = "dataitem"
)

var (
	API_SERVER         = Env("API_SERVER", false)
	API_PORT           = Env("API_PORT", false)
	USER_TP_ADMIN  int = 2
	USER_TP_UNKNOW     = -1
)

type Quota struct {
	Rep_Private int
	Rep_Public  int
}

func auth(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	if login_Name == "" {
		http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
	}
	c.Map(login_Name)
	return

}

func getUserType(r *http.Request, db *DB) int {
	login_Name := r.Header.Get("User")
	token := r.Header.Get(AUTHORIZATION)
	if login_Name == "" || token == "" {
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

func updateUser(loginName, token string, real interface{}) ([]byte, error) {
	if loginName == "" || token == "" {
		Log.Infof("create repository token: %s login_name: %s", token, loginName)
		return nil, nil
	}

	b, _ := json.Marshal(real)

	return HttpPostJson(fmt.Sprintf("http://%s:%s/quota/%s/repository/use", API_SERVER, API_PORT, loginName), b, AUTHORIZATION, token)
}

func getQuota(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	loginName := r.Header.Get("User")
	token := r.Header.Get(AUTHORIZATION)
	if loginName == "" || token == "" {
		Log.Infof("create repository token: %s, login_name: %s", token, loginName)
		http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
		return
	}

	q := getUserQuota(token, loginName)
	c.Map(q)
}

func authAdmin(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {

	if getUserType(r, db) != USER_TP_ADMIN {
		Log.Infof("auth admin: %s", r.Header.Get("User"))
		http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
		return
	}
	login_Name := r.Header.Get("User")
	c.Map(login_Name)
}

func chkRepPermission(w http.ResponseWriter, r *http.Request, param martini.Params, c martini.Context, db *DB) {
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
		return
	}
	repName := strings.TrimSpace(param["repname"])
	if repName == "" {
		http.Error(w, ErrNoParameter("repname").ErrToString(), 401)
		return
	}
	delUser := strings.TrimSpace(param["username"])

	rep, _ := db.getRepository(bson.M{COL_REPNAME: repName})
	if rep.Create_user != user {
		http.Error(w, E(ErrorCodePermissionDenied).ErrToString(), 401)
		return
	}
	c.Map(Rep_Permission{Repository_name: repName, User_name: delUser})
	c.Map(rep.Repaccesstype)
}

func chkItemPermission(w http.ResponseWriter, r *http.Request, param martini.Params, c martini.Context, db *DB) {
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
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
		http.Error(w, E(ErrorCodePermissionDenied).ErrToString(), 401)
		return
	}
	c.Map(Item_Permission{Repository_name: repName, Dataitem_name: itemname})
}
