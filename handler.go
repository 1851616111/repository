package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ACCESS_PRIVATE   = "private"
	ACCESS_PUBLIC    = "public"
	COL_REP_NAME     = "repository_name"
	COL_REP_ACC      = "repaccesstype"
	COL_ITEM_NAME    = "dataitem_name"
	COL_TAG_TAG      = "tag"
	COL_SELECT_LABEL = "labelname"
	COL_PERMIT_USER  = "user_name"
	PAGE_INDEX       = 1
	PAGE_SIZE        = 3
)

//curl http://127.0.0.1:8080/repositories/rep12asda232312sd -d "{\"repaccesstype\": \"public\",\"comment\": \"中国移动北京终端详情\",
//\"label\":{\"sys\":{\"name\":\"中国移动\"},\"opt\":{\"name\":\"中国移动\"},\"owner\":{\"name\":\"中国移动\"},\"other\":{\"name\":\"中国移动\"}}}" -H admin:admin
func createRHandler(r *http.Request, rsp *Rsp, param martini.Params, login_name string) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	body, _ := ioutil.ReadAll(r.Body)

	rep := new(repository)
	if len(body) > 0 {
		if err := json.Unmarshal(body, &rep); err != nil {
			return rsp.Json(400, ErrParseJson(err))
		}
	}

	now := time.Now()
	if rep.Repaccesstype != ACCESS_PUBLIC || rep.Repaccesstype != ACCESS_PRIVATE {
		rep.Repaccesstype = ACCESS_PUBLIC
	}
	rep.Optime = now
	rep.Ct = now
	rep.Create_user = login_name
	rep.Repository_name = repname
	rep.Stars, rep.Items = 0, 0

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

//curl http://127.0.0.1:8080/repositories/rep123 -X DELETE -H admin:admin
func delRHandler(r *http.Request, rsp *Rsp, param martini.Params, loginName string) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{"repository_name": repname}
	rep, err := db.getRepository(Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	if rep.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}
	if err := db.delRepository(Q); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories
//curl http://10.1.235.98:8080/repositories
//curl http://10.1.235.98:8080/repositories?page=2&size=3
func getRsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, userName string) (int, string) {
	//		page_index, page_size := PAGE_INDEX, PAGE_SIZE
	//		if p := strings.TrimSpace(r.FormValue("page")); p != "" {
	//			if page_index, _ = strconv.Atoi(p); page_index <= 0 {
	//				return rsp.Json(400, ErrInvalidParameter("page"))
	//			}
	//
	//		}
	//		if p := strings.TrimSpace(r.FormValue("size")); p != "" {
	//			if page_size, _ = strconv.Atoi(p); page_size <= 0 {
	//				return rsp.Json(400, ErrInvalidParameter("size"))
	//			}
	//		}
	////		var Q bson.M
	//		if p := strings.TrimSpace(r.FormValue("username")); p != "" {
	//			Q = bson.M{C_REPOSITORY_PERMIT: ACCESS_PRIVATE,}
	//		}
	//
	l := []dataItem{}
	//	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(nil).Sort("ct").Skip((PAGE_INDEX - 1) * PAGE_SIZE).Limit(PAGE_SIZE).All(&l); err != nil {
	//		rsp.Json(400, ErrDataBase(err))
	//	}

	return rsp.Json(200, E(OK), l)
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23 -d "{\"itemaccesstype\":\"public\", \"meta\":\"{}\",\"sample\":\"{}\",\"comment\":\"中国移动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\",\"refresh\":\"3天\"}}}" -H admin:admin
func createDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	Q := bson.M{COL_REP_NAME: repname}
	if _, err := db.getRepository(Q); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("repname : %s", repname)))
	}

	body, _ := ioutil.ReadAll(r.Body)

	d := new(dataItem)
	if len(body) > 0 {
		if err := json.Unmarshal(body, &d); err != nil {
			return rsp.Json(400, ErrParseJson(err))
		}
	}

	d.Repository_name = repname
	d.Dataitem_name = itemname
	d.Create_name = loginName
	now := time.Now()
	d.Optime = now
	d.Ct = now
	d.Stars, d.Tags = 0, 0

	if d.Itemaccesstype != ACCESS_PRIVATE || d.Itemaccesstype != ACCESS_PUBLIC {
		d.Itemaccesstype = ACCESS_PUBLIC
	}

	if err := db.DB(DB_NAME).C(C_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, ErrDataBase(err))

	}
	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23 -d "{\"repaccesstype\":\"public\", \"meta\":\"{}\",\"sample\":\"{}\",\"comment\":\"中国移动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\",\"refresh\":\"3天\"}}}" -H user:admin
func updateDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	body, _ := ioutil.ReadAll(r.Body)

	d := new(dataItem)
	if len(body) > 0 {
		if err := json.Unmarshal(body, &d); err != nil {
			return rsp.Json(400, ErrParseJson(err))
		}
	}
//	if
//
//	d.Repository_name = repname
//	d.Dataitem_name = itemname
//	d.Create_name = loginName
//	now := time.Now()
//	d.Optime = now
//	d.Ct = now
//	d.Stars, d.Tags = 0, 0

	if d.Itemaccesstype != ACCESS_PRIVATE || d.Itemaccesstype != ACCESS_PUBLIC {
		d.Itemaccesstype = ACCESS_PUBLIC
	}

	if err := db.DB(DB_NAME).C(C_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, ErrDataBase(err))

	}
	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/rep123/bear23 -X DELETE -H admin:admin
func delDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)

	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("repname : %s", repname)))
	}
	if item.Create_name != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	if err := db.delDataitem(Q); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:8080/select_labels/CHINA -d "order=100" -H admin:admin
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

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -d "{\"comment\":\"this is a tag\"}" -H user:admin
func setTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}
	tagname := param["tag"]
	if tagname == "" {
		return rsp.Json(400, ErrNoParameter("tag"))
	}

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)

	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("itemname : %s", itemname)))
	}

	if item.Create_name != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	t := new(tag)
	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		if err := json.Unmarshal(body, &t); err != nil {
			return rsp.Json(400, ErrParseJson(err))
		}
	}
	t.Repository_name, t.Dataitem_name, t.Tag = repname, itemname, tagname
	t.Optime = time.Now()

	if err := db.DB(DB_NAME).C(C_TAG).Insert(t); err != nil {
		return rsp.Json(400, ErrDataBase(err))

	}
	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -H user:admin
func getTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}
	tagname := param["tag"]
	if tagname == "" {
		return rsp.Json(400, ErrNoParameter("tag"))
	}

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname, COL_TAG_TAG: tagname}
	tag, err := db.getTag(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("tag : %s", tag)))
	}

	return rsp.Json(200, E(OK), tag)
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -H user:admin -X DELETE
func delTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}
	tagname := param["tag"]
	if tagname == "" {
		return rsp.Json(400, ErrNoParameter("tag"))
	}

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname, COL_TAG_TAG: tagname}
	err := db.delTag(Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23
func getDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK), item)
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

//curl http://127.0.0.1:8080/permit/michael -H michael:pan
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

//curl http://127.0.0.1:8080/permit/michael -d "repname=rep00002" -H user:admin
func setUsrPmtRepsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	var user_name, repname string
	if user_name = strings.TrimSpace(param["user_name"]); user_name == "" {
		return rsp.Json(400, ErrNoParameter("user_name"))
	}
	if repname = strings.TrimSpace(r.PostFormValue("repname")); repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{COL_REP_NAME: repname, COL_REP_ACC: ACCESS_PRIVATE}

	if _, err := db.getRepository(Q); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	Exec := bson.M{"repository_name": repname, "user_name": user_name}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMIT).Insert(Exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}
