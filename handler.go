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
	REPACCESS_PRIVATE        = "private"
	REPACCESS_PUBLIC         = "public"
	COLUMN_USER_NAME         = "LOGIN_NAME"
	COLUMN_REP_NAME          = "repository_name"
	COLUMN_ITEM_NAME         = "dataitem_name"
	PORTAL_REQUEST_TP_CHOSEN = "chosen"
	PAGE_INDEX               = 1
	PAGE_SIZE                = 3
)

//curl http://10.1.235.98:6000/repositories/rep123 -d ""
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

	if err := db.DB(DB_NAME_MONGO).C(M_REPOSITORY).Insert(rep); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:6000/repositories/rep00001
func getRHandler(r *http.Request, rsp *Rsp, param martini.Params) (int, string) {
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	res := new(repository)
	Q := bson.M{"repository_name": repname}
	if err := db.DB(DB_NAME_MONGO).C(M_REPOSITORY).Find(Q).One(res); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK), res)
}

//curl http://10.1.235.98:6000/repositories
//curl http://10.1.235.98:6000/repositories?page=2&size=3
func getRsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	//		page_index, page_size := PAGE_INDEX, PAGE_SIZE
	//		if p := strings.TrimSpace(r.FormValue("page")); p != "" {
	//			if page_index, _ = strconv.Atoi(p); page_index <= 0 {
	//				return rsp.Json(400, "page can not little than 0")
	//			}
	//
	//		}
	//		if p := strings.TrimSpace(r.FormValue("size")); p != "" {
	//			if page_size, _ = strconv.Atoi(p); page_size <= 0 {
	//				return rsp.Json(400, "size can not little than 0")
	//			}
	//		}
	//
	//		res := []Result{}
	//
	//		db.DB(DB_NAME).C(M_DATAITEM).Find().Sort("ct").Skip((PAGE_INDEX -1) *  PAGE_SIZE).Limit(PAGE_SIZE)
	////
	////		l, err := db.getDataitems("", "", (page_index-1)*page_size, page_size)
	//		if err != nil {
	//			return rsp.Json(400, err.Error())
	//		}
	//		m := make(M)
	//
	//		for _, v := range l {
	//			if s, exists := m[v.Repository_name]; exists {
	//				s = append(s.([]DataItem), v)
	//				m[v.Repository_name] = s
	//			} else {
	//				m[v.Repository_name] = []DataItem{v}
	//			}
	//		}
	//		for k, v := range m {
	//			res = append(res, Result{k.(string), REPACCESS_PUBLIC, v.([]DataItem)})
	//		}

	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:6000/repositories/NBA/bear -d ""
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

	if err := db.DB(DB_NAME_MONGO).C(M_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, ErrDataBase(err))

	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:6000/select_labels/CHINA -d "order=100"
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
	if err := db.DB(DB_NAME_MONGO).C(M_SELECT).Insert(s); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

////curl http://10.1.235.98:6000/selects?select_labels=NBA
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

//curl http://10.1.235.98:6000/select_labels
func getSelectLabelsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	l := []Select{}
	err := db.DB(DB_NAME_MONGO).C(M_SELECT).Find(nil).Sort("-order").All(&l)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK) ,l)
}



//curl http://10.1.235.98:6000/repositories/NBA/bear/0001
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

//curl http://10.1.235.98:6000/repositories/位置信息大全/全国在网（新增）终端
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

	Q := bson.M{COLUMN_REP_NAME: repname, COLUMN_ITEM_NAME: itemname}
	if err := db.DB(DB_NAME_MONGO).C(M_DATAITEM).Find(Q).One(d); err != nil {
		return *d, err
	}
	return *d, nil
}

//curl http://10.1.235.98:6000/selects -d "repname=NBA&itemname=bear&select_labels=h"
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
	if _, err := db.DB(DB_NAME_MONGO).C(M_DATAITEM).Upsert(Q, m); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}

//curl http://10.1.235.98:6000/selects?select_labels=h
func getSelectsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {

	var m bson.M

	if select_labels := strings.TrimSpace(r.FormValue("select_labels")); select_labels != "" {
		m = bson.M{"label.sys.select_labels": select_labels}
	} else {
		m = bson.M{"label.sys.select_labels": bson.M{"$exists": true}}
	}

	l := []names{}
	if err := db.DB(DB_NAME_MONGO).C(M_DATAITEM).Find(m).Sort("-label.sys.order").All(&l); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK), l)
}
