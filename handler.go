package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"strings"
	"time"
)

const (
	REPACCESS_PRIVATE        = "private"
	REPACCESS_PUBLIC         = "public"
	USER_NAME                = "LOGIN_NAME"
	REP_NAME                 = "REPOSITORY_NAME"
	PORTAL_REQUEST_TP_CHOSEN = "chosen"
)

//curl http://127.0.0.1:8088/repositories
func getRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	res := []Result{}

	l, err := db.getDataitems("", "")
	if err != nil {
		return rsp.Json(400, err.Error())
	}
	m := make(M)

	for _, v := range l {
		if s, exists := m[v.Repository_name]; exists {
			s = append(s.([]DataItem), v)
			m[v.Repository_name] = s
		} else {
			m[v.Repository_name] = []DataItem{v}
		}
	}
	for k, v := range m {
		res = append(res, Result{k.(string), REPACCESS_PUBLIC, v.([]DataItem)})
	}

	return rsp.Json(200, res)
}

//curl http://127.0.0.1:8088/repositories/NBA/bear -d "" -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
func setDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, user_id string) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, "no param repname")
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, "no param itemname")
	}

	if l := db.getRepository(REP_NAME, repname); len(l) == 0 {
		return rsp.Json(400, "repname do not exist")
	}
	d := new(DataItem)

	d.ParseRequeset(r)
	d.BuildRequeset(repname, itemname, user_id)
	if err := db.setDataitem(d); err != nil {
		return rsp.Json(400, err.Error())
	}

	return rsp.Json(200, "ok")
}

//curl http://127.0.0.1:8088/repository/cba/items
func getRepoByNameHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	res, err := db.getDataitems("repository_name", repname)
	if err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, res)
}

//curl http://127.0.0.1:8088/repositories/chosen -d "chosen_name=NBA&dataitem_id=1011"
func setItemChoseHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	d := new(Dataitem_Chosen)
	if err := d.ParseRequeset(r); err != nil {
		return rsp.Json(400, err.Error())
	}
	if err := db.setDataitem_Chosen(d); err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, "ok")
}

//curl http://127.0.0.1:8088/repositories/chosen/dataitem?chosen_name=NBA
//curl http://127.0.0.1:8088/repositories/chosen?chosen_name=cb
func getItemChoseHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	var l []Dataitem_Chosen
	var err error

	if chosen_name := strings.TrimSpace(r.FormValue("chosen_name")); chosen_name != "" {
		l, err = db.getDataitem_Chosen(chosen_name)
	} else {
		l, err = db.getDataitem_Chosen()
	}

	if err != nil {
		return rsp.Json(400, err.Error())
	}

	return rsp.Json(200, l)
}

//curl http://54.223.244.55:8088/repositories/chosen/names
func getChosenNamesHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	l, err := db.getDataitem_ChosenNames()
	if err != nil {
		return rsp.Json(400, err)
	}
	return rsp.Json(200, l)
}

//curl http://10.1.51.32:8080/subscriptions/login -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
func login(r *http.Request, rsp *Rsp) (int, string) {
	return 200, "ok"
}

//curl http://127.0.0.1:8088/repositories/chosen
func getItemsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	var err error
	l := []Data{}
	var l_s []Dataitem_Chosen
	if chosen_name := strings.TrimSpace(r.FormValue("chosen_name")); chosen_name != "" {
		l_s, err = db.getDataitem_Chosen(chosen_name)
	} else {
		l_s, err = db.getDataitem_Chosen()
	}
	get(err)

	l_str := []interface{}{}
	for _, v := range l_s {
		l_str = append(l_str, v.Dataitem_id)
	}
	l = db.getDataitemsByIds(l_str)

	return rsp.Json(200, l)
}

//curl http://127.0.0.1:8088/repositories/NBA/bear/0001 -d "" -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
func setTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	d, err := getDMid(r, rsp, param, db)
	if d == nil {
		return rsp.Json(400, "no found dataitem")
	}
	get(err)

	tag := param["tag"]
	if tag == "" {
		return rsp.Json(400, "no param tag")
	}

	t := new(Tag)
	t.ParseRequeset(r)
	t.Tag = tag
	t.Optime = time.Now().Format(TimeFormat)
	t.Dataitem_id = d.Dataitem_id

	if err := db.setTag(t); err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, "ok")
}

//curl http://127.0.0.1:8080/subscriptions/NBA/bear
//curl http://127.0.0.1:8088/subscriptions/NBA/bear
//curl http://127.0.0.1:8088/inner/NBA/bear
//curl http://127.0.0.1:8088/inner/NBA/bear/tags
//curl http://127.0.0.1:8088/repositories/位置信息大全/全国在网（新增）终端
func getDataitemHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	d, err := getDMid(r, rsp, param, db)
	if d == nil {
		return rsp.Json(400, "no found dataitem")
	}
	get(err)

	res := Data{Item: d}

//	l := []Tag{}

//	m, err := db.getDataitemUsageByIds(d.Dataitem_id)
//	if err != nil {
//		return rsp.Json(400, err.Error())
//	}
//	l, err = db.getTags(d.Dataitem_id)
//	if err != nil {
//		return rsp.Json(400, err.Error())
//	}
//	dd := m[d.Dataitem_id]
//	res.Usage = &dd
//	res.Tags = l

	//		l, err = db.getTags(d.Dataitem_id)
	//		if err != nil {
	//			return rsp.Json(400, err.Error())
	//		}
	//		res.Tags = l

	return rsp.Json(200, res)
}
