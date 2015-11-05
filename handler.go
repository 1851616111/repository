package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"strings"
)

const (
	REPACCESS_PRIVATE        = "private"
	REPACCESS_PUBLIC         = "public"
	USER_NAME                = "USER_ID"
	REP_NAME                 = "REPOSITORY_NAME"
	PORTAL_REQUEST_TP_CHOSEN = "chosen"
)

//curl http://127.0.0.1:8080/repositories -u panxy3@asiainfo.com:88888888
func getRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, user_id string) (int, string) {
	res := []Result{}

	l, err := db.getDataitem(USER_NAME, user_id)
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

//curl http://127.0.0.1:8080/repositories/NBA/bear -d "" -u panxy3@asiainfo.com:88888888
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

//curl http://127.0.0.1:8080/repository/cba/items
func getRepoByNameHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	res, err := db.getDataitem("repository_name", repname)
	if err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, res)
}

//curl http://127.0.0.1:8080/dataitem/chosen -d "chosen_name=NBA&dataitem_id=1000001"
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

//curl http://127.0.0.1:8080/portal/portal/dataitem/chosen?chosen_name=cb
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

//curl http://127.0.0.1:8080/portal/dataitem/chosen/names
func getChosenNamesHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	l, err := db.getDataitem_ChosenNames()
	if err != nil {
		return rsp.Json(400, err)
	}
	return rsp.Json(200, l)
}

//curl http://10.1.51.32:8080/subscriptions/login -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
func login(r *http.Request, rsp *Rsp) (int, string) {
	return rsp.Json(200, "ok")
}

//curl http://10.1.51.32:8080/portal/dataitem?tp=chosen
func getItemsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	tp := r.FormValue("tp")
	var err error
	l := []Data{}
	switch tp {
	case PORTAL_REQUEST_TP_CHOSEN:
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

		usages_m, err := db.getDataitemUsageByIds(l_str...)
		get(err)
		items_m, err := db.getDataitemByIds(l_str...)
		get(err)

		m := make(M)
		for k, item := range items_m {
			usage := usages_m[k]
			usage.Refresh_date = buildTime(usage.Refresh_date)
			d := Data{item, usage}
			m[k] = d
		}

		for _, v := range m {
			l = append(l, v.(Data))
		}

	}

	return rsp.Json(200, l)
}

//curl http://10.1.51.32:8080/portal/dataitem/test
//func test(r *http.Request, rsp *Rsp, db *DB) (int, string) {
//
////	usages_m, err  := db.getDataitemUsageByIds()
//	items_m, err := db.getDataitemByIds()
//	get(err)
////	log.Println(usages_m)
//	log.Println(items_m)
//	return rsp.Json(200, "ok")
//}
