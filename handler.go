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
	COL_REPNAME         = "repository_name"
	COL_REP_ACC         = "repaccesstype"
	COL_REP_ITEMS       = "items"
	COL_ITEM_NAME       = "dataitem_name"
	COL_ITEM_ACC        = "itemaccesstype"
	COL_ITEM_TAGS       = "tags"
	COL_COMMENT         = "comment"
	COL_PRICE           = "price"
	COL_CREATE_USER     = "create_user"
	COL_LABEL           = "label"
	COL_OPTIME          = "optime"
	COL_ITEM_META       = "meta"
	COL_ITEM_SAMPLE     = "sample"
	COL_TAG_NAME        = "tag"
	COL_SELECT_LABEL    = "labelname"
	COL_SELECT_ORDER    = "order"
	COL_SELECT_ICON     = "icon"
	COL_PERMIT_USER     = "user_name"
	COL_PERMIT_REPNAME  = "repository_name"
	COL_PERMIT_ITEMNAME = "dataitem_name"
	COL_PERMIT_WRITE    = "write"
	PAGE_INDEX          = 1
	PAGE_SIZE           = 6
	PAGE_SIZE_SEARCH    = 10
	PAGE_SIZE_SELECT    = 10
	LABEL_NED_CHECK     = "supply_style"
	SUPPLY_STYLE_API    = "api"
	SUPPLY_STYLE_BATCH  = "batch"
	SUPPLY_STYLE_FLOW   = "flow"
	CMD_INC             = "$inc"
	CMD_SET             = "$set"
	CMD_UNSET           = "$unset"
	CMD_IN              = "$in"
	CMD_OR              = "$or"
	PREFIX_META         = "meta"
	PREFIX_SAMPLE       = "sample"
	MQ_TYPE_ADD_TAG     = "0x00020000"
	MQ_TYPE_DEL_TAG     = "0x00020001"
	MQ_TYPE_DEL_ITEM    = "0x00020002"
	MQ_TYPE_DEL_REP     = "0x00020003"
)

var (
	SUPPLY_STYLE_ALL = []string{SUPPLY_STYLE_API, SUPPLY_STYLE_BATCH, SUPPLY_STYLE_FLOW}
	NED_CHECK_LABELS = []string{LABEL_NED_CHECK}
)

func createRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, login_name string, l Limit) (int, string) {
	defer db.Close()
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err:", err)
	}

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

	var opt interface{}
	has := db.countNum(C_REPOSITORY, bson.M{COL_CREATE_USER: login_name, COL_REP_ACC: rep.Repaccesstype})
	max := 0
	switch rep.Repaccesstype {
	case ACCESS_PUBLIC:
		max = l.Rep_Public
		opt = struct {
			Public int `json:"public"`
		}{has}
	case ACCESS_PRIVATE:
		max = l.Rep_Private
		opt = struct {
			Private int `json:"private"`
		}{has}
	}

	if has >= max {
		return rsp.Json(400, ErrRepOutOfLimit(max))
	}

	rep.Optime = now.String()
	rep.Ct = now
	rep.Create_user = login_name
	rep.Repository_name = repname
	rep.Items = 0
	rep.chkLabel()

	if err := db.DB(DB_NAME).C(C_REPOSITORY).Insert(rep); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	if _, err := updateUser(r, opt); err != nil {
		Log.Errorf("create dataitem update User err: %s", err.Error())
	}
	return rsp.Json(200, E(OK))
}

func getRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}
	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size < -1 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}

	user := r.Header.Get("User")
	showItems := strings.TrimSpace(r.FormValue("items"))

	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REPNAME, repname)))
	}
	rep.Optime = buildTime(rep.Optime)

	if rep.Repaccesstype == ACCESS_PRIVATE {
		Q = bson.M{COL_PERMIT_REPNAME: repname, COL_PERMIT_USER: user}
		if user != "" {
			if rep.Create_user != user && !db.hasPermission(C_REPOSITORY_PERMISSION, Q) {
				Log.Debugf("[Auth] login name %s, repository name %s.", user, rep.Create_user)
				return rsp.Json(400, E(ErrorCodePermissionDenied))
			}
		} else {
			return rsp.Json(400, E(ErrorCodePermissionDenied))
		}
	}

	var res struct {
		repository
		Dataitems []string `json:"dataitems,omitempty"`
	}
	res.repository = rep

	if showItems != "" {
		items := []string{}
		ds := []dataItem{}
		Q := bson.M{COL_REPNAME: repname, COL_CREATE_USER: user}
		ds, err = db.getDataitems(page_index, page_size, Q)
		get(err)
		for _, v := range ds {
			items = append(items, v.Dataitem_name)
		}
		res.Dataitems = items
	}

	return rsp.Json(200, E(OK), res)
}

//curl http://127.0.0.1:8080/repositories/rep123 -X DELETE -H admin:admin
func delRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REPNAME, repname)))
	}
	if rep.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}
	if err := db.delRepository(Q); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}

	tmp := m_rep{Type: MQ_TYPE_DEL_REP, Repository_name: Q[COL_REPNAME], Time: time.Now().String()}
	go func(rep m_rep) {
		msg.MqJson(rep)
	}(tmp)

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA -d "{\"repaccesstype\":\"public\",\"comment\":\"中国移动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\",\"refresh\":\"3天\"}}}" -H user:admin -X PUT
func updateRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	Q := bson.M{COL_REPNAME: repname}
	repo, err := db.getRepository(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REPNAME, repname)))
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

	selector := bson.M{COL_REPNAME: repname}

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

	if len(u) > 0 {
		u[COL_OPTIME] = time.Now().String()
		updater := bson.M{"$set": u}
		go asynUpdateOpt(C_REPOSITORY, selector, updater)
	}

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8089/repositories
func getRsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}

	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size < -1 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}

	targetName := strings.TrimSpace(r.FormValue("username"))
	loginName := r.Header.Get("User")
	Q := bson.M{}
	if loginName != "" && targetName == "" { //login already and search myrepositories
		Q = bson.M{COL_CREATE_USER: loginName}
	} else if loginName == "" && targetName != "" { // no login nd search targetName
		Q = bson.M{COL_CREATE_USER: targetName, COL_REP_ACC: ACCESS_PUBLIC}
	} else if loginName != "" && targetName != "" { //loging and search targetName
		q := []bson.M{}
		p_reps, _ := db.getPermits(COL_PERMIT_REPNAME, bson.M{COL_PERMIT_USER: loginName})
		l := []string{}
		if l_p_reps, ok := p_reps.([]Rep_Permission); ok {
			if len(l_p_reps) > 0 {
				for _, v := range p_reps.([]Rep_Permission) {
					l = append(l, v.Repository_name)
				}
				q_private := bson.M{}
				q_private[COL_REPNAME] = bson.M{CMD_IN: l}
				q_private[COL_CREATE_USER] = targetName
				q_private[COL_REP_ACC] = ACCESS_PRIVATE
				q = append(q, q_private)
			}
		}

		q_public := bson.M{COL_CREATE_USER: targetName, COL_REP_ACC: ACCESS_PUBLIC}

		q = append(q, q_public)
		switch len(q) {
		case 1:
			Q = q_public
		case 2:
			Q[CMD_OR] = q
		}
	} else if loginName == targetName && targetName == "" { // invalid param
		return rsp.Json(400, ErrNoParameter("username or user"))
	}

	rep := []repository{}

	if page_size == -1 {
		if err := db.DB(DB_NAME).C(C_REPOSITORY).Find(Q).Sort("-ct").Select(bson.M{COL_REPNAME: "1"}).All(&rep); err != nil {
			return rsp.Json(400, ErrDataBase(err))
		}
	} else {
		err := db.DB(DB_NAME).C(C_REPOSITORY).Find(Q).Sort("-ct").Select(bson.M{COL_REPNAME: "1"}).Skip((page_index - 1) * page_size).Limit(page_size).All(&rep)
		if err != nil {
			return rsp.Json(400, ErrDataBase(err))
		}
	}

	l := []names{}
	for _, v := range rep {
		l = append(l, names{Repository_name: v.Repository_name})
	}

	return rsp.Json(200, E(OK), l)
}

func createDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	Q := bson.M{COL_REPNAME: repname}
	if _, err := db.getRepository(Q); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("repname : %s", repname)))
	}

	if n, _ := db.DB(DB_NAME).C(C_DATAITEM).Find(bson.M{COL_REPNAME: repname}).Count(); n >= 30 {
		return rsp.Json(400, E(ErrorCodeItemOutOfLimit))
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err ", err)

	}
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	d := new(dataItem)
	if err := json.Unmarshal(body, &d); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	d.Repository_name = repname
	d.Dataitem_name = itemname
	d.Create_user = loginName
	now := time.Now()
	d.Optime = now.String()
	d.Ct = now
	d.Tags = 0

	if d.Itemaccesstype != ACCESS_PRIVATE && d.Itemaccesstype != ACCESS_PUBLIC {
		d.Itemaccesstype = ACCESS_PUBLIC
	}

	if d.Label == "" {
		return rsp.Json(400, ErrNoParameter("label"))
	}

	if err := ifInLabel(d.Label, LABEL_NED_CHECK); err != nil {
		return rsp.Json(400, err)
	}

	if d.Price == "" {
		return rsp.Json(400, ErrNoParameter("price"))
	}

	if err := chkPrice(d.Price, getSupplyStyleTp(d.Label)); err != nil {
		return rsp.Json(400, err)
	}

	if d.Meta != "" {
		if err := db.setFile(PREFIX_META, repname, itemname, []byte(d.Meta)); err != nil {
			return rsp.Json(400, err)
		}
	}

	if d.Sample != "" {
		if err := db.setFile(PREFIX_SAMPLE, repname, itemname, []byte(d.Sample)); err != nil {
			return rsp.Json(400, err)
		}
	}

	if err := db.DB(DB_NAME).C(C_DATAITEM).Insert(d); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	go asynUpdateOpt(C_REPOSITORY, bson.M{COL_REPNAME: repname}, bson.M{CMD_INC: bson.M{"items": 1}, CMD_SET: bson.M{COL_OPTIME: now.String()}})

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23 -d "{\"itemaccesstype\":\"public\", \"meta\":\"{}\",\"sample\":\"{}\",\"comment\":\"中国移动北京终端详情\", \"label\":{\"sys\":{\"supply_style\":\"flow\",\"refresh\":\"3天\"}}}" -H user:admin
func updateDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s %s:=%s", COL_REPNAME, repname, COL_ITEM_NAME, itemname)))
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

	selector := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	u := bson.M{}

	if d.Itemaccesstype != "" {
		if d.Itemaccesstype != ACCESS_PRIVATE && d.Itemaccesstype != ACCESS_PUBLIC {
			return rsp.Json(400, ErrInvalidParameter("itemaccesstype"))
		}
		u[COL_ITEM_ACC] = d.Itemaccesstype
	}

	if d.Price != nil {
		if err := chkPrice(d.Price, getSupplyStyleTp(d.Label)); err != nil {
			return rsp.Json(400, err)
		}
		u[COL_PRICE] = d.Price
	}

	if d.Dataitem_name != "" {
		u[COL_ITEM_NAME] = d.Dataitem_name
	}

	if d.Meta != "" {
		if err := db.setFile(PREFIX_META, repname, itemname, []byte(d.Meta)); err != nil {
			return rsp.Json(400, err)
		}
	}

	if d.Sample != "" {
		if err := db.setFile(PREFIX_SAMPLE, repname, itemname, []byte(d.Sample)); err != nil {
			return rsp.Json(400, err)
		}
	}

	if d.Comment != "" {
		u[COL_COMMENT] = d.Comment
	}

	if len(u) > 0 {
		now := time.Now().String()
		u[COL_OPTIME] = now
		updater := bson.M{"$set": u}
		go asynUpdateOpt(C_DATAITEM, selector, updater)

		go asynUpdateOpt(C_REPOSITORY, bson.M{COL_REPNAME: repname}, bson.M{"$set": bson.M{COL_OPTIME: now}})
	}
	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/rep123/bear23 -X DELETE -H admin:admin
func delDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}

	item, err := db.getDataitem(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s %s:=%s", COL_REPNAME, repname, COL_ITEM_NAME, itemname)))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	if err := db.delDataitem(Q); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}

	go func(db *DB) {
		db.delFile(PREFIX_META, repname, itemname)
		db.delFile(PREFIX_SAMPLE, repname, itemname)
	}(db.copy())

	go asynUpdateOpt(C_REPOSITORY, bson.M{COL_REPNAME: repname}, bson.M{CMD_INC: bson.M{"items": -1}, CMD_SET: bson.M{COL_OPTIME: time.Now().String()}})

	tmp := m_item{Type: MQ_TYPE_DEL_ITEM, Repository_name: Q[COL_REPNAME], Dataitem_name: Q[COL_ITEM_NAME], Time: time.Now().String()}
	go func(item m_item) {
		msg.MqJson(tmp)
	}(tmp)

	return rsp.Json(200, E(OK))
}

func setSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	labelname := ""
	if labelname = strings.TrimSpace(param["labelname"]); labelname == "" {
		return rsp.Json(400, ErrNoParameter("labelname"))
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err :", err)
	}

	s := new(Select)
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	if err := json.Unmarshal(body, &s); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}
	s.LabelName = labelname

	if err := db.DB(DB_NAME).C(C_SELECT).Insert(s); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}

func updateSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	labelname := ""
	if labelname = strings.TrimSpace(param["labelname"]); labelname == "" {
		return rsp.Json(400, ErrNoParameter("labelname"))
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err", err)
	}
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
	if s.Icon != "" {
		u[COL_SELECT_ICON] = s.Icon
	}

	if len(u) == 0 {
		return rsp.Json(400, ErrNoParameter("newlabelname or order"))
	}

	updater := bson.M{CMD_SET: u}
	go asynUpdateOpt(C_SELECT, selector, updater)

	return rsp.Json(200, E(OK))
}

func delSelectLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
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

func getSelectLabelsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	defer db.Close()

	l, ll := []Select{Select{LabelName: "全部精选"}}, []Select{}
	err := db.DB(DB_NAME).C(C_SELECT).Find(nil).Sort("-order").All(&ll)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	l = append(l, ll...)
	return rsp.Json(200, E(OK), l)
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -d "{\"comment\":\"this is a tag\"}" -H user:admin
func createTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
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

	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("itemname : %s", itemname)))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read request body err", err)
	}
	if len(body) == 0 {
		return rsp.Json(400, ErrNoParameter(""))
	}
	t := new(tag)
	if err := json.Unmarshal(body, &t); err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	t.Repository_name, t.Dataitem_name, t.Tag = repname, itemname, tagname
	now := time.Now().String()
	t.Optime = now

	if err := db.DB(DB_NAME).C(C_TAG).Insert(t); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	go asynUpdateOpt(C_REPOSITORY, bson.M{COL_REPNAME: repname}, bson.M{CMD_SET: bson.M{COL_OPTIME: now}})

	go asynUpdateOpt(C_DATAITEM, Q, bson.M{CMD_INC: bson.M{"tags": 1}, CMD_SET: bson.M{COL_OPTIME: now}})

	go func(t tag) {
		m_t := m_tag{Type: MQ_TYPE_ADD_TAG, Repository_name: t.Repository_name, Dataitem_name: t.Dataitem_name, Tag: t.Tag, Time: t.Optime}
		msg.MqJson(m_t)
	}(*t)

	return rsp.Json(200, E(OK))
}

func updateTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
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

	Q_item := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q_item)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("itemname : %s", itemname)))
	}

	Q_tag := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname, COL_TAG_NAME: tagname}
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
	go asynUpdateOpt(C_REPOSITORY, bson.M{COL_REPNAME: repname}, bson.M{CMD_SET: bson.M{COL_OPTIME: now}})

	go asynUpdateOpt(C_DATAITEM, Q_item, bson.M{CMD_INC: bson.M{"tags": 1}, CMD_SET: bson.M{COL_OPTIME: now}})

	go asynUpdateOpt(C_TAG, Q_tag, bson.M{"$set": bson.M{COL_COMMENT: t.Comment, COL_OPTIME: now}})

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -H user:admin
func getTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
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

	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname, COL_TAG_NAME: tagname}
	tag, err := db.getTag(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("tag : %s", tagname)))
	}
	tag.Optime = buildTime(tag.Optime)

	return rsp.Json(200, E(OK), tag)
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -H user:admin -X DELETE
func delTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
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

	var item dataItem
	var err error
	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	if item, err = db.getDataitem(Q); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf("itemname : %s", itemname)))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	Q[COL_TAG_NAME] = tagname
	err = db.delTag(Q)
	if err != nil && err != mgo.ErrNotFound {
		return rsp.Json(400, ErrDataBase(err))
	}

	go asynUpdateOpt(C_DATAITEM, bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}, bson.M{CMD_INC: bson.M{"tags": -1}, CMD_SET: bson.M{COL_OPTIME: time.Now().String()}})

	t := m_tag{Type: MQ_TYPE_DEL_TAG, Repository_name: Q[COL_REPNAME], Dataitem_name: Q[COL_ITEM_NAME], Tag: Q[COL_TAG_NAME], Time: time.Now().String()}
	go func(t m_tag) {
		msg.MqJson(t)
	}(t)

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8089/repositories/mobile/app
func getDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}
	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size < -1 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}

	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REPNAME, repname)))
	}

	user := r.Header.Get("User")
	if rep.Repaccesstype == ACCESS_PRIVATE {
		Q := bson.M{COL_PERMIT_REPNAME: rep.Repository_name, COL_PERMIT_USER: user}
		if user != "" {
			if rep.Create_user != user && !db.hasPermission(COL_PERMIT_REPNAME, Q) {
				return rsp.Json(400, E(ErrorCodePermissionDenied))
			}
		} else {
			return rsp.Json(400, E(ErrorCodePermissionDenied))
		}
	}

	Q = bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s,%s=%s ", COL_REPNAME, repname, COL_ITEM_NAME, itemname)))
	}
	item.Optime = buildTime(item.Optime)

	b_m, err := db.getFile(PREFIX_META, repname, itemname)
	if err != nil {
		Log.Error("get dataitem meta err :", err)
		item.Meta = ""
	} else {
		item.Meta = strings.TrimSpace(string(b_m))
	}

	b_s, err := db.getFile(PREFIX_SAMPLE, repname, itemname)
	if err != nil {
		Log.Error("get dataitem sample err :", err)
		item.Sample = ""
	} else {
		item.Sample = strings.TrimSpace(string(b_s))
	}

	tags, err := db.getTags(page_index, page_size, Q)
	get(err)
	buildTagsTime(tags)

	var res struct {
		dataItem
		Tags []tag `json:"taglist"`
	}
	res.dataItem = item
	res.Tags = tags

	return rsp.Json(200, E(OK), res)
}

//curl http://127.0.0.1:8089/repositories/mobile/app/subpermission
func getDWithPermissionHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	user := r.Header.Get("User")
	if user == "" {
		return rsp.Json(400, E(ErrorCodeUnauthorized))
	}

	var err error
	var rep repository
	var item dataItem
	Q := bson.M{COL_REPNAME: repname}
	rep, err = db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s", COL_REPNAME, repname)))
	}

	Q = bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err = db.getDataitem(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s=%s,%s=%s ", COL_REPNAME, repname, COL_ITEM_NAME, itemname)))
	}

	Q = bson.M{COL_PERMIT_REPNAME: repname, COL_PERMIT_USER: user}
	hasRepPermission := db.hasPermission(C_REPOSITORY_PERMISSION, Q)
	Q = bson.M{COL_PERMIT_ITEMNAME: itemname, COL_PERMIT_USER: user}
	hasItemPermission := db.hasPermission(C_DATAITEM_PERMISSION, Q)

	if rep.Repaccesstype == ACCESS_PRIVATE && !hasRepPermission {
		return rsp.Json(400, E(ErrorCodePermissionDenied), false)
	}

	if item.Itemaccesstype == ACCESS_PRIVATE && !hasItemPermission {
		return rsp.Json(400, E(ErrorCodePermissionDenied), false)
	}

	return rsp.Json(200, E(OK), item)
}

func updateSelectHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	repname := strings.TrimSpace(param["repname"])
	itemname := strings.TrimSpace(param["itemname"])
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

	selector := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}

	if n, _ := db.DB(DB_NAME).C(C_DATAITEM).Find(selector).Count(); n == 0 {
		return rsp.Json(400, ErrQueryNotFound(fmt.Sprintf(" %s %s", repname, itemname)))
	}

	u := bson.M{}
	u["label.sys.select_labels"] = select_labels
	u["label.sys.order"] = order
	updater := bson.M{"$set": u}

	go q_c.producer(exec{C_DATAITEM, selector, updater})

	return rsp.Json(200, E(OK))
}

func delSelectHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	repname := strings.TrimSpace(param["repname"])
	itemname := strings.TrimSpace(param["itemname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	selector := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}

	u := bson.M{}
	u["label.sys.select_labels"] = 1
	u["label.sys.order"] = 1
	updater := bson.M{CMD_UNSET: u}

	go q_c.producer(exec{C_DATAITEM, selector, updater})

	return rsp.Json(200, E(OK))
}

func getSelectsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	defer db.Close()

	page_index, page_size := PAGE_INDEX, PAGE_SIZE_SEARCH
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}

	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size <= 0 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}

	username := r.Header.Get("User")
	var m bson.M
	if select_labels := strings.TrimSpace(r.FormValue("select_labels")); select_labels != "" {
		m = bson.M{"label.sys.select_labels": select_labels}
	} else {
		m = bson.M{"label.sys.select_labels": bson.M{"$exists": true}}
	}

	Q := bson.M{}
	l := db.getPublicReps()
	if username != "" {
		private := db.getPrivateReps(username)
		l = append(l, private...)
	}

	if len(l) > 0 {
		q := bson.M{COL_REPNAME: bson.M{CMD_IN: l}}
		Q[CMD_AND] = []bson.M{q, m}
	} else {
		Q = m
	}

	res := []names{}
	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(Q).Sort("-label.sys.order").Skip((page_index - 1) * page_size).Limit(page_size).All(&res); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	total, _ := db.DB(DB_NAME).C(C_DATAITEM).Find(Q).Count()
	result := struct {
		Namelist `json:"select"`
		Total    int `json:"total"`
	}{
		res,
		total,
	}

	return rsp.Json(200, E(OK), result)
}

//curl http://127.0.0.1:8080/permit/michael -H michael:pan
func getUsrPmtRepsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	var user_name string
	if user_name = strings.TrimSpace(param["user_name"]); user_name == "" {
		return rsp.Json(400, ErrNoParameter("user_name"))
	}

	l := []Rep_Permission{}
	Q := bson.M{"user_name": user_name}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Find(Q).All(&l); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK), l)
}

//curl http://127.0.0.1:8080/permit/michael -d "repname=rep00002" -H user:admin
func setUsrPmtRepsHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	defer db.Close()
	var user_name, repname string
	if user_name = strings.TrimSpace(param["user_name"]); user_name == "" {
		return rsp.Json(400, ErrNoParameter("user_name"))
	}
	if repname = strings.TrimSpace(r.PostFormValue("repname")); repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}
	Q := bson.M{COL_REPNAME: repname, COL_REP_ACC: ACCESS_PRIVATE}

	if _, err := db.getRepository(Q); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	Exec := bson.M{COL_REPNAME: repname, "user_name": user_name}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMISSION).Insert(Exec); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	return rsp.Json(200, E(OK))
}
