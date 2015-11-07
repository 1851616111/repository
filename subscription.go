package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//curl http://127.0.0.1:8080/subscriptions -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
//curl http://10.1.235.96:8080/subscriptions -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
func getSHandler(r *http.Request, rsp *Rsp, db *DB, user_id string) (int, string) {
	url := fmt.Sprintf("http://10.1.235.96:8081/subscriptions/user/%s", user_id)
	b, err := HttpGet(url)
	get(err)
	s := struct {
		Err  string         `json:"error,omitempty"`
		Subs []Subscription `json:"subscriptions"`
	}{}
	err = json.Unmarshal(b, &s)
	get(err)
	l := []interface{}{}
	for _, v := range s.Subs {
		l = append(l, v.Dataitem_id)
	}
	d := db.getDataitemsByIds(l)
	return rsp.Json(200, d)
}

func setSHandler(r *http.Request, rsp *Rsp, db *DB, user_id string) (int, string) {

	return rsp.Json(200, "ok")
}
