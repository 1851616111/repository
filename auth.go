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
	VIP_SERVICE_RET_RET_PUB   = "repoPub"
	VIP_SERVICE_RET_RET_PRI   = "repoPri"
	VIP_SERVICE_ADMIN_PUB     = -1
	VIP_SERVICE_ADMIN_PRI     = -1
)

var (
	API_SERVER         = Env("API_SERVER", false)
	API_PORT           = Env("API_PORT", false)
	USER_TP_ADMIN  int = 2
	USER_TP_UNKNOW     = -1
)

type Limit struct {
	Rep_Private int `json:"private"`
	Rep_Public  int `json:"public"`
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

func updateUser(r *http.Request, l Limit) ([]byte, error) {
	login_Name := r.Header.Get("User")
	token := r.Header.Get(AUTHORIZATION)
	if login_Name == "" || token == "" {
		Log.Infof("create repository token: %sm login_name: %s", token, login_Name)
		return nil, nil
	}
	b, _ := json.Marshal(l)

	return HttpPostJson(fmt.Sprintf("http://%s:%s/quota/%s/repository/use", API_SERVER, API_PORT, login_Name), string(b), AUTHORIZATION, token)
}

func chkUserLimit(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	token := r.Header.Get(AUTHORIZATION)
	//	token := "Token 1ad5d08d974ae394a511af1c0a0a7d79"
	if login_Name == "" || token == "" {
		Log.Infof("create repository token: %sm login_name: %s", token, login_Name)
		http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
		return
	}
	b, err := httpGet(fmt.Sprintf("http://%s:%s/vip/%s", API_SERVER, API_PORT, login_Name), AUTHORIZATION, token)
	get(err)
	result := new(Result)
	err = json.Unmarshal(b, result)
	get(err)
	if result.Data != nil {
		u := result.Data.(map[string]interface{})
		l := Limit{}
		if pub, exist := u[VIP_SERVICE_RET_RET_PUB]; exist {
			l.Rep_Public = int(pub.(float64))
		}
		if pri, exist := u[VIP_SERVICE_RET_RET_PRI]; exist {
			l.Rep_Private = int(pri.(float64))
		}
		if l.Rep_Public == VIP_SERVICE_ADMIN_PUB {
			l.Rep_Public = 100000
		}
		if l.Rep_Private == VIP_SERVICE_ADMIN_PUB {
			l.Rep_Private = 100000
		}
		c.Map(l)
		return
	}
	http.Error(w, E(ErrorCodeUnauthorized).ErrToString(), 401)
}

func authAdmin(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {

	if getUserType(r, db) != USER_TP_ADMIN {
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

	if rep, _ := db.getRepository(bson.M{COL_REPNAME: repName}); rep.Create_user != user {
		http.Error(w, E(ErrorCodePermissionDenied).ErrToString(), 401)
		return
	}
	c.Map(Rep_Permission{Repository_name: repName})
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
