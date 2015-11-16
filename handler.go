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
	ACCESS_PRIVATE      = "private"
	ACCESS_PUBLIC       = "public"
	COL_REP_NAME        = "repository_name"
	COL_REP_ACC         = "repaccesstype"
	COL_ITEM_NAME       = "dataitem_name"
	COL_ITEM_ACC        = "itemaccesstype"
	COL_COMMENT         = "comment"
	COL_LABEL           = "label"
	COL_OPTIME          = "optime"
	COL_ITEM_META       = "meta"
	COL_ITEM_SAMPLE     = "sample"
	COL_TAG_NAME        = "tag"
	COL_SELECT_LABEL    = "labelname"
	COL_SELECT_ORDER    = "order"
	COL_PERMIT_USER     = "user_name"
	PAGE_INDEX          = 1
	PAGE_SIZE           = 3
	LABEL_NED_CHECK     = "supply_style"
	SUPPLY_STYLE_SINGLE = "single"
	SUPPLY_STYLE_BATCH  = "batch"
	SUPPLY_STYLE_FLOW   = "flow"
	CMD_INC             = "$inc"
	CMD_SET             = "$set"
)

var (
	SUPPLY_STYLE_ALL = []string{SUPPLY_STYLE_SINGLE, SUPPLY_STYLE_BATCH, SUPPLY_STYLE_FLOW}
	NED_CHECK_LABELS = []string{LABEL_NED_CHECK}
)

func createRHandler(r *http.Request, rsp *Rsp, param martini.Params, login_name string) (int, string) {

	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	body, _ := ioutil.ReadAll(r.Body)
	rep := new(repository)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &rep); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}
	now := time.Now()
	if rep.Repaccesstype != ACCESS_PUBLIC && rep.Repaccesstype != ACCESS_PRIVATE {
		rep.Repaccesstype = ACCESS_PUBLIC
	}
	rep.Optime = now.String()
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
	Q := bson.M{COL_REP_NAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REP_NAME, repname)))
	}
	rep.Optime = buildTime(rep.Optime)
	return rsp.Json(200, E(OK), rep)
}

//curl http://127.0.0.1:8080/repositories/rep123 -X DELETE -H admin:admin
func delRHandler(r *http.Request, rsp *Rsp, param martini.Params, loginName string) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{COL_REP_NAME: repname}
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

//curl http://127.0.0.1:8080/repositories/NBA -d "{\"repaccesstype\":\"public\",\"comment\":\"中国移动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\",\"refresh\":\"3天\"}}}" -H user:admin -X PUT
func updateRHandler(r *http.Request, rsp *Rsp, param martini.Params, loginName string) (int, string) {
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	Q := bson.M{COL_REP_NAME: repname}
	repo, err := db.getRepository(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REP_NAME, repname)))
	}

	if repo.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	rep := new(repository)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &rep); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	selector := bson.M{COL_REP_NAME: repname}

	u := bson.M{}
	if rep.Repaccesstype != "" {
		if rep.Repaccesstype != ACCESS_PRIVATE && rep.Repaccesstype != ACCESS_PUBLIC {
			return rsp.Json(400, ErrInvalidParameter("repaccesstype"))
		}
		u[COL_REP_ACC] = rep.Repaccesstype
	}

	if rep.Comment != "" {
		u[COL_COMMENT] = rep.Comment
	}

	if rep.Label != "" {
		u[COL_LABEL] = rep.Label
	}

	if len(u) > 0 {
		u[COL_OPTIME] = time.Now().String()
		updater := bson.M{"$set": u}
		go asynOpt(C_REPOSITORY, selector, updater)
	}

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories
//curl http://10.1.235.98:8080/repositories
//curl http://10.1.235.98:8080/repositories?page=2&size=3
func getRsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
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

//curl http://127.0.0.1:8080/repositories/NBA/bear23 -d "{\"itemaccesstype\":\"public\", \"meta\":\"{}\",\"sample\":\"{}\",\"comment\":\"中国移 动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\"}}}" -H user:admin
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
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &d); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	d.Repository_name = repname
	d.Dataitem_name = itemname
	d.Create_user = loginName
	now := time.Now()
	d.Optime = now.String()
	d.Ct = now
	d.Stars, d.Tags = 0, 0

	if d.Itemaccesstype != ACCESS_PRIVATE && d.Itemaccesstype != ACCESS_PUBLIC {
		d.Itemaccesstype = ACCESS_PUBLIC
	}

	if err := ifInLabel(d.Label, LABEL_NED_CHECK); err != nil {
		return rsp.Json(400, err)
	}

	if err := db.DB(DB_NAME).C(C_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	go asynOpt(C_REPOSITORY, bson.M{COL_REP_NAME: repname}, bson.M{CMD_INC: bson.M{"items": 1}, CMD_SET: bson.M{COL_OPTIME: now.String()}})

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23 -d "{\"itemaccesstype\":\"public\", \"meta\":\"{}\",\"sample\":\"{}\",\"comment\":\"中国移动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\",\"refresh\":\"3天\"}}}" -H user:admin
func updateDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
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
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s %s:=%s", COL_REP_NAME, repname, COL_ITEM_NAME, itemname)))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	d := new(dataItem)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &d); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	selector := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	u := bson.M{}

	if d.Itemaccesstype != "" {
		if d.Itemaccesstype != ACCESS_PRIVATE && d.Itemaccesstype != ACCESS_PUBLIC {
			return rsp.Json(400, ErrInvalidParameter("itemaccesstype"))
		}
		u[COL_ITEM_ACC] = d.Itemaccesstype
	}

	if d.Dataitem_name != "" {
		u[COL_ITEM_NAME] = d.Dataitem_name
	}

	if d.Meta != "" {
		u[COL_ITEM_META] = d.Meta
	}

	if d.Sample != "" {
		u[COL_ITEM_SAMPLE] = d.Sample
	}

	if d.Comment != "" {
		u[COL_COMMENT] = d.Comment
	}

	if d.Label != "" {
		if err := ifInLabel(d.Label, LABEL_NED_CHECK); err != nil {
			return rsp.Json(400, err)
		}
		u[COL_LABEL] = d.Label
	}

	if len(u) > 0 {
		now := time.Now().String()
		u[COL_OPTIME] = now
		updater := bson.M{"$set": u}
		go asynOpt(C_DATAITEM, selector, updater)

		go asynOpt(C_REPOSITORY, bson.M{COL_REP_NAME: repname}, bson.M{"$set": bson.M{COL_OPTIME: now}})
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

	Q := bson.M{COL_REP_NAME: repname}
	go asynOpt(C_REPOSITORY, Q, bson.M{CMD_INC: bson.M{"items": -1}, CMD_SET: bson.M{COL_OPTIME: time.Now().String()}})

	Q[COL_ITEM_NAME] = itemname
	item, err := db.getDataitem(Q)

	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s %s:=%s", COL_REP_NAME, repname, COL_ITEM_NAME, itemname)))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	if err := db.delDataitem(Q); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}

func setSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	labelname := ""
	if labelname = strings.TrimSpace(param["labelname"]); labelname == "" {
		return rsp.Json(400, ErrNoParameter("labelname"))
	}

	s := new(Select)
	s.LabelName = labelname

	if err := db.DB(DB_NAME).C(C_SELECT).Insert(s); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

func updateSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	labelname := ""
	if labelname = strings.TrimSpace(param["labelname"]); labelname == "" {
		return rsp.Json(400, ErrNoParameter("labelname"))
	}

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}

	s := new(Select)
	if err := json.Unmarshal(body, &s); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	selector := bson.M{COL_SELECT_LABEL: labelname}
	if _, err := db.getSelect(selector); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("labelname : %s", labelname)))
	}

	u := bson.M{}
	if s.NewLabelName != "" {
		u[COL_SELECT_LABEL] = s.NewLabelName
	}
	if s.Order > 0 {
		u[COL_SELECT_ORDER] = s.Order
	}

	if len(u) == 0 {
		return rsp.Json(400, ErrNoParameter("newlabelname or order"))
	}

	updater := bson.M{CMD_SET: u}
	go asynOpt(C_SELECT, selector, updater)
	return rsp.Json(200, E(OK))
}

func delSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	labelname := ""
	if labelname = strings.TrimSpace(param["labelname"]); labelname == "" {
		return rsp.Json(400, ErrNoParameter("labelname"))
	}
	Q := bson.M{COL_SELECT_LABEL: labelname}
	if err := db.delSelect(Q); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("labelname : %s", labelname)))
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

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	now := time.Now().String()
	t := new(tag)
	body, _ := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &t); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	t.Repository_name, t.Dataitem_name, t.Tag = repname, itemname, tagname
	t.Optime = now

	go asynOpt(C_REPOSITORY, bson.M{COL_REP_NAME: repname}, bson.M{CMD_SET: bson.M{COL_OPTIME: now}})

	go asynOpt(C_DATAITEM, Q, bson.M{CMD_INC: bson.M{"tags": 1}, CMD_SET: bson.M{COL_OPTIME: now}})

	if err := db.DB(DB_NAME).C(C_TAG).Insert(t); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}

func updateTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
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

	Q_item := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q_item)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("itemname : %s", itemname)))
	}

	Q_tag := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname, COL_TAG_NAME: tagname}
	if _, err := db.getTag(Q_tag); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("tagname : %s", tagname)))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	t := new(tag)
	body, _ := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &t); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	if t.Comment == "" {
		return rsp.Json(400, E(ErrorCodeInvalidParameters))
	}

	now := time.Now().String()
	go asynOpt(C_REPOSITORY, bson.M{COL_REP_NAME: repname}, bson.M{CMD_SET: bson.M{COL_OPTIME: now}})

	go asynOpt(C_DATAITEM, Q_item, bson.M{CMD_INC: bson.M{"tags": 1}, CMD_SET: bson.M{COL_OPTIME: now}})

	go asynOpt(C_TAG, Q_tag, bson.M{"$set": bson.M{COL_COMMENT: t.Comment, COL_OPTIME: now}})

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

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname, COL_TAG_NAME: tagname}
	tag, err := db.getTag(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("tag : %s", tagname)))
	}
	tag.Optime = buildTime(tag.Optime)
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

	Q := bson.M{COL_REP_NAME: repname, COL_ITEM_NAME: itemname}
	if _, err := db.getDataitem(Q); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("itemname : %s", itemname)))
	}

	go asynOpt(C_DATAITEM, Q, bson.M{CMD_INC: bson.M{"tags": -1}, CMD_SET: bson.M{COL_OPTIME: time.Now().String()}})

	Q[COL_TAG_NAME] = tagname
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
	item.Optime = buildTime(item.Optime)
	return rsp.Json(200, E(OK), item)
}

//curl http://10.1.235.98:8080/selects -d "repname=NBA&itemname=bear&select_labels=h"
func updateLabelHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	repname := strings.TrimSpace(r.FormValue("repname"))
	itemname := strings.TrimSpace(r.FormValue("itemname"))
	select_labels := strings.TrimSpace(r.FormValue("select_labels"))
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	if select_labels == "" {
		return rsp.Json(400, ErrNoParameter("select_labels"))
	}
	order := 1
	if o := strings.TrimSpace(r.FormValue("order")); o != "" {
		order, _ = strconv.Atoi(o)
	}

	selector := bson.M{COL_REP_NAME: repname}
	collectionName := C_REPOSITORY

	if itemname != "" {
		collectionName = C_DATAITEM
		selector[COL_ITEM_NAME] = itemname
	}

	if n, _ := db.DB(DB_NAME).C(collectionName).Find(selector).Count(); n == 0 {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s %s", repname, itemname)))
	}

	u := bson.M{}
	u["label.sys.select_labels"] = select_labels
	u["order"] = order
	updater := bson.M{"$set": u}

	go q_c.producer(exec{collectionName, selector, updater})

	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:8080/selects -d "repname=NBA&itemname=bear&select_labels=h" -X DELETE
func deleteLabelHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	repname := strings.TrimSpace(r.FormValue("repname"))
	itemname := strings.TrimSpace(r.FormValue("itemname"))
	select_labels := strings.TrimSpace(r.FormValue("select_labels"))
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	if select_labels == "" {
		return rsp.Json(400, ErrNoParameter("select_labels"))
	}
	order := 1
	if o := strings.TrimSpace(r.FormValue("order")); o != "" {
		order, _ = strconv.Atoi(o)
	}

	selector := bson.M{COL_REP_NAME: repname}
	collectionName := C_REPOSITORY

	if itemname != "" {
		collectionName = C_DATAITEM
		selector[COL_ITEM_NAME] = itemname
	}

	if n, _ := db.DB(DB_NAME).C(collectionName).Find(selector).Count(); n == 0 {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s %s", repname, itemname)))
	}

	u := bson.M{}
	u["label.sys.select_labels"] = select_labels
	u["order"] = order
	updater := bson.M{"$set": u}

	go q_c.producer(exec{collectionName, selector, updater})

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
	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(m).Limit(PAGE_SIZE).Sort("-label.sys.order").All(&l); err != nil {
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

	Exec := bson.M{COL_REP_NAME: repname, "user_name": user_name}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMIT).Insert(Exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}
