package main

import (
	"errors"
	"github.com/go-martini/martini"
	"github.com/lunny/log"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ACCESS_PRIVATE   = "private"
	ACCESS_PUBLIC    = "public"
	COL_REP_NAME     = "repository_name"
	COL_ITEM_NAME    = "dataitem_name"
	COL_REP_ACC      = "repaccesstype"
	COL_SELECT_LABEL = "labelname"
	COL_PERMIT_USER  = "user_name"
	PAGE_INDEX       = 1
	PAGE_SIZE        = 3
)

//curl http://10.1.235.98:8080/repositories/rep123 -d ""
//func createRHandler(r *http.Request, rsp *Rsp, param martini.Params, login_name string) (int, string) {
func createRHandler(r *http.Request, rsp *Rsp, param martini.Params) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	rep := new(repository)
	rep.ParseRequeset(r)
	rep.BuildRequest()
	rep.Create_user = "panxy3@asiainfo.com"
	rep.Repository_name = repname
	rep.Ct = time.Now()

	if err := db.DB(DB_NAME).C(C_REPOSITORY).Insert(rep); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:8080/repositories/rep00001
func getRHandler(r *http.Request, rsp *Rsp, param martini.Params) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{"repository_name": repname}
	rep, err := db.getRepository(Q)
	return rsp.Json(200, ErrDataBase(err), rep)
}

//curl http://127.0.0.1:8080/repositories
//curl http://10.1.235.98:8080/repositories
//curl http://10.1.235.98:8080/repositories?page=2&size=3
func getRsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	//	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	//	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
	//		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
	//			return rsp.Json(400, ErrInvalidParameter("page"))
	//		}
	//
	//	}
	//	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
	//		if page_size, _ = strconv.Atoi(p); page_size <= 0 {
	//			return rsp.Json(400, ErrInvalidParameter("size"))
	//		}
	//	}
	//	var Q bson.M
	//	if p := strings.TrimSpace(r.FormValue("username")); p != "" {
	//		Q = bson.M{C_REPOSITORY_PERMIT:}
	//	}

	l := []dataItem{}
	//	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(nil).Sort("ct").Skip((PAGE_INDEX - 1) * PAGE_SIZE).Limit(PAGE_SIZE).All(&l); err != nil {
	//		rsp.Json(400, ErrDataBase(err))
	//	}

	return rsp.Json(200, E(OK), l)
}

//curl http://10.1.235.98:8080/repositories/NBA/bear -d ""
func createDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	//		if l := db.getRepository(COLUMN_REP_NAME, repname); len(l) == 0 {
	//			return rsp.Json(400, "repname do not exist")
	//		}
	d := new(dataItem)

	d.ParseRequeset(r)
	d.BuildRequeset(repname, itemname, "panxy3@asiainfo.com")

	log.Printf("%+v", d)

	if err := db.DB(DB_NAME).C(C_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, ErrDataBase(err))

	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:8080/select_labels/CHINA -d "order=100"
func setSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	labelname := ""
	if labelname = strings.TrimSpace(param["labelname"]); labelname == "" {
		return rsp.Json(400, ErrNoParameter("labelname"))
	}
	s := Select{LabelName: labelname}

	if order := strings.TrimSpace(r.FormValue("order")); order != "" {
		o, err := strconv.Atoi(order)
		if err != nil {
			return rsp.Json(400, ErrInvalidParameter("order"))
		}
		s.Order = o
	}
	if err := db.DB(DB_NAME).C(C_SELECT).Insert(s); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

////curl http://10.1.235.98:8080/selects?select_labels=NBA
//func getSelectsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
//		var l []Select
//		var err error
//		Q := bson.M{}
//		if select_labels := strings.TrimSpace(r.FormValue("select_labels")); select_labels != "" {
//			Q["select_labels"] = select_labels
//		}
//		if err := db.DB(DB_NAME).C(M_SELECT).Insert(s); err != nil {
//			return rsp.Json(400, err.Error())
//		}
//		if err != nil {
//			return rsp.Json(400, err.Error())
//		}
//
//	return rsp.Json(200, OK)
//}

//curl http://10.1.235.98:8080/select_labels
func getSelectLabelsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	l := []Select{}
	err := db.DB(DB_NAME).C(C_SELECT).Find(nil).Sort("-order").All(&l)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), l)
}

//curl http://10.1.235.98:8080/repositories/NBA/bear/0001
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
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:8080/repositories/位置信息大全/全国在网（新增）终端
func getDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	//
	//	d, err := getDMid(r, rsp, param, db)
	//	if err != nil {
	//		return rsp.Json(400, "get dataitem err"+err.Error())
	//	}
	//
	//	//	l, err := db.getTags(d.Dataitem_id)
	//	//	if err != nil {
	//	//		return rsp.Json(400, err.Error())
	//	//	}
	//
	//	//	res := Data{Item: d, Tags: l}
	return rsp.Json(200, E(OK))
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

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(Q).One(d); err != nil {
		return *d, err
	}
	return *d, nil
}

//curl http://10.1.235.98:8080/selects -d "repname=NBA&itemname=bear&select_labels=h"
func updateLabelHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {

	var m bson.M
	repname := strings.TrimSpace(r.FormValue("repname"))
	itemname := strings.TrimSpace(r.FormValue("itemname"))
	select_labels := strings.TrimSpace(r.FormValue("select_labels"))
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}
	if select_labels == "" {
		return rsp.Json(400, ErrNoParameter("select_labels"))
	}
	order := 1
	if o := strings.TrimSpace(r.FormValue("order")); o != "" {
		order, _ = strconv.Atoi(o)
	}

	mm := bson.M{"select_labels": select_labels, "order": order}
	m = bson.M{"$set": bson.M{"label.sys": mm}}

	// check if dataitem exists

	Q := bson.M{"repository_name": repname, "dataitem_name": itemname}
	if _, err := db.DB(DB_NAME).C(C_DATAITEM).Upsert(Q, m); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:8080/selects?select_labels=终端信息
func getSelectsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	var m bson.M

	if select_labels := strings.TrimSpace(r.FormValue("select_labels")); select_labels != "" {
		m = bson.M{"label.sys.select_labels": select_labels}
	} else {
		m = bson.M{"label.sys.select_labels": bson.M{"$exists": true}}
	}

	l := []names{}
	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(m).Sort("-label.sys.order").Limit(PAGE_SIZE).All(&l); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK), l)
}

//curl http://127.0.0.1:8080/permit/michael -u michael:pan
func getUsrPmtRepsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	var user_name string
	if user_name = strings.TrimSpace(param["user_name"]); user_name == "" {
		return rsp.Json(400, ErrNoParameter("user_name"))
	}

	l := []Repository_Permit{}
	Q := bson.M{"user_name": user_name}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMIT).Find(Q).All(&l); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), l)
}

//curl http://127.0.0.1:8080/permit/michael -d "repname=rep00002" -u michael:pan
func setUsrPmtRepsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	var user_name, repname string
	if user_name = strings.TrimSpace(param["user_name"]); user_name == "" {
		return rsp.Json(400, ErrNoParameter("user_name"))
	}
	if repname = strings.TrimSpace(r.PostFormValue("repname")); repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{COL_REP_NAME: repname, COL_REP_ACC: ACCESS_PRIVATE}
	log.Println("-------->", Q)
	if _, err := db.getRepository(Q); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	Exec := bson.M{"repository_name": repname, "user_name": user_name}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMIT).Insert(Exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}
