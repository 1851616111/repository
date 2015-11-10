package main

import (
	"errors"
	"github.com/asiainfoLDP/datahub_repository/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/go-martini/martini"
	"github.com/lunny/log"
	"net/http"
	"strings"
)

const (
	REPACCESS_PRIVATE        = "private"
	REPACCESS_PUBLIC         = "public"
	COLUMN_USER_NAME         = "LOGIN_NAME"
	COLUMN_REP_NAME          = "repository_name"
	COLUMN_ITEM_NAME         = "dataitem_name"
	PORTAL_REQUEST_TP_CHOSEN = "chosen"
	PAGE_INDEX               = 1
	PAGE_SIZE                = 5
	OK                       = "OK"
)

//curl http://127.0.0.1:8088/repositories/rep123 -d ""
//func createRHandler(r *http.Request, rsp *Rsp, param martini.Params, login_name string) (int, string) {
func createRHandler(r *http.Request, rsp *Rsp, param martini.Params) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, "no param repname")
	}
	rep := new(repository)
	rep.ParseRequeset(r)
	rep.BuildRequest()
	rep.Create_user = "panxy3@asiainfo.com"
	rep.Repository_name = repname

	if err := db.DB(DB_NAME).C(M_REPOSITORY).Insert(rep); err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, OK)
}

//curl http://127.0.0.1:8088/repositories/rep123
func getRHandler(r *http.Request, rsp *Rsp, param martini.Params) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, "no param repname")
	}
	res := new(repository)
	Q := bson.M{"repository_name": repname}
	if err := db.DB(DB_NAME).C(M_REPOSITORY).Find(Q).One(res); err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, res)
}

//curl http://127.0.0.1:8088/repositories
//curl http://127.0.0.1:8088/repositories?page=2&size=3
func getRsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, login_name string) (int, string) {
	//	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	//	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
	//		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
	//			return rsp.Json(400, "page can not little than 0")
	//		}
	//
	//	}
	//	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
	//		if page_size, _ = strconv.Atoi(p); page_size <= 0 {
	//			return rsp.Json(400, "size can not little than 0")
	//		}
	//	}
	//
	//	res := []Result{}
	//
	//	l, err := db.getDataitems("", "", (page_index-1)*page_size, page_size)
	//	if err != nil {
	//		return rsp.Json(400, err.Error())
	//	}
	//	m := make(M)
	//
	//	for _, v := range l {
	//		if s, exists := m[v.Repository_name]; exists {
	//			s = append(s.([]DataItem), v)
	//			m[v.Repository_name] = s
	//		} else {
	//			m[v.Repository_name] = []DataItem{v}
	//		}
	//	}
	//	for k, v := range m {
	//		res = append(res, Result{k.(string), REPACCESS_PUBLIC, v.([]DataItem)})
	//	}

	return rsp.Json(200, OK)
}

//curl http://127.0.0.1:8088/repositories/NBA/bear -d ""
func createDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, "no param repname")
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, "no param itemname")
	}

	//		if l := db.getRepository(COLUMN_REP_NAME, repname); len(l) == 0 {
	//			return rsp.Json(400, "repname do not exist")
	//		}
	d := new(dataItem)

	d.ParseRequeset(r)
	d.BuildRequeset(repname, itemname, "panxy3@asiainfo.com")

	log.Printf("%+v", d)

	if err := db.DB(DB_NAME).C(M_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, "ok")
}

//curl http://127.0.0.1:8088/select_labels/CHINA -d "dataitem_name=NBA&repository_name=bear&labelname=basketball"
func setSelectLabelHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	labelname := ""
	if labelname = strings.TrimSpace(r.FormValue("labelname")); labelname == "" {
		return rsp.Json(400, "no param label")
	}
	d := new(Select)
	if err := d.ParseRequeset(r); err != nil {
		return rsp.Json(400, err.Error())
	}
	d.LabelName = labelname

	if err := m_db.DB(DB_NAME).C("select").Insert(d); err != nil {
		return rsp.Json(400, err.Error())
	}
	return rsp.Json(200, "ok")
}

//curl http://127.0.0.1:8088/selects?select_labels=NBA
func getSelectsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	//	var l []Select
	//	var err error
	//
	//	if select_labels := strings.TrimSpace(r.FormValue("select_labels")); select_labels != "" {
	//		l, err = db.getSelects(select_labels)
	//	} else {
	//		l, err = db.getSelects()
	//	}
	//
	//	if err != nil {
	//		return rsp.Json(400, err.Error())
	//	}

	return rsp.Json(200, OK)
}

//curl http://54.223.244.55:8088/select_labels
func getSelectLabelsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	//	l, err := db.getSelectLabels()
	//	if err != nil {
	//		return rsp.Json(400, err)
	//	}
	return rsp.Json(200, OK)
}

//curl http://10.1.51.32:8080/subscriptions/login -u panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909
func login(r *http.Request, rsp *Rsp) (int, string) {
	return 200, "ok"
}

//curl http://127.0.0.1:8088/repositories/chosen
func getItemsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	//	var err error
	//	l := []Data{}
	//	var l_s []Select
	//	if chosen_name := strings.TrimSpace(r.FormValue("chosen_name")); chosen_name != "" {
	//		l_s, err = db.getSelects(chosen_name)
	//	} else {
	//		l_s, err = db.getSelects()
	//	}
	//	get(err)
	//
	//	l_str := []interface{}{}
	//	for _, v := range l_s {
	//		l_str = append(l_str, v.Dataitem_name)
	//	}
	//	l = db.getDataitemsByIds(l_str)

	return rsp.Json(200, OK)
}

//curl http://127.0.0.1:8088/repositories/NBA/bear/0001
func setTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	//	d, err := getDMid(r, rsp, param, db)
	//	if d == nil {
	//		return rsp.Json(400, "no found dataitem")
	//	}
	//	get(err)
	//
	//	tag := param["tag"]
	//	if tag == "" {
	//		return rsp.Json(400, "no param tag")
	//	}
	//
	//	t := new(Tag)
	//	t.ParseRequeset(r)
	//	t.Tag = tag
	//	t.Optime = time.Now().Format(TimeFormat)
	//	t.Dataitem_id = d.Dataitem_id
	//
	//	if err := db.setTag(t); err != nil {
	//		return rsp.Json(400, err.Error())
	//	}
	return rsp.Json(200, "ok")
}

//curl http://127.0.0.1:8088/repositories/位置信息大全/全国在网（新增）终端
func getDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {

	d, err := getDMid(r, rsp, param, db)
	if err != nil {
		return rsp.Json(400, "get dataitem err" + err.Error())
	}


//	l, err := db.getTags(d.Dataitem_id)
//	if err != nil {
//		return rsp.Json(400, err.Error())
//	}

//	res := Data{Item: d, Tags: l}
	return rsp.Json(200, d)
}
func getDMid(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (dataItem, error) {
	d := new(dataItem)
	repname := param["repname"]
	if repname == "" {
		return *d, errors.New("no param repname")
	}
	itemname := param["itemname"]
	if itemname == "" {
		return *d, errors.New("no param repname")
	}

	Q := bson.M{COLUMN_REP_NAME: repname, COLUMN_ITEM_NAME: itemname}
	if err := db.DB(DB_NAME).C(M_DATAITEM).Find(Q).One(d); err != nil {
		return *d, err
	}
	return *d, nil
}
