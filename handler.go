package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ACCESS_PRIVATE            = "private"
	ACCESS_PUBLIC             = "public"
	COL_REPNAME               = "repository_name"
	COL_REP_ACC               = "repaccesstype"
	COL_REP_ITEMS             = "items"
	COL_ITEM_NAME             = "dataitem_name"
	COL_ITEM_ACC              = "itemaccesstype"
	COL_ITEM_TAGS             = "tags"
	COL_COMMENT               = "comment"
	COL_RANK                  = "rank"
	COL_PRICE                 = "price"
	COL_CREATE_USER           = "create_user"
	COL_REP_COOPERATOR        = "cooperators"
	COL_ITEM_COOPERATOR       = "ifcooperate"
	COL_LABEL                 = "label"
	COL_OPTIME                = "optime"
	COL_ITEM_META             = "meta"
	COL_ITEM_SAMPLE           = "sample"
	COL_TAG_NAME              = "tag"
	COL_SELECT_LABEL          = "labelname"
	COL_SELECT_ORDER          = "order"
	COL_SELECT_ICON           = "icon"
	COL_PERMIT_USER           = "user_name"
	COL_PERMIT_REPNAME        = "repository_name"
	COL_PERMIT_ITEMNAME       = "dataitem_name"
	COL_PERMIT_WRITE          = "write"
	PAGE_INDEX                = 1
	PAGE_SIZE                 = 6
	PAGE_SIZE_SEARCH          = 10
	PAGE_SIZE_SELECT          = 10
	LIMIT_TAG_LENGTH          = 100
	LIMIT_ITEM_LENGTH         = 100
	LIMIT_REP_LENGTH          = 52
	LIMIT_COMMENT_LENGTH      = 600
	PARAM_TAG_NAME            = "tag"
	PARAM_ITEM_NAME           = "itemname"
	PARAM_REP_NAME            = "repname"
	PARAM_COMMENT_NAME        = "comment"
	COL_ITEM_SYPPLY_STYLE     = "supply_style"
	LABEL_NED_CHECK           = COL_ITEM_SYPPLY_STYLE
	SUPPLY_STYLE_API          = "api"
	SUPPLY_STYLE_BATCH        = "batch"
	SUPPLY_STYLE_FLOW         = "flow"
	SUPPLY_STYLE_UNRECOGNIZED = ""
	CMD_INC                   = "$inc"
	CMD_ADDTOSET              = "$addToSet"
	CMD_SET                   = "$set"
	CMD_UNSET                 = "$unset"
	CMD_IN                    = "$in"
	CMD_OR                    = "$or"
	CMD_REGEX                 = "$regex"
	CMD_OPTION                = "$options"
	CMD_AND                   = "$and"
	CMD_CASE_ALL              = "$i"
	CMD_PULL                  = "$pull"
	CMD_NOTEQUAL              = "$ne"
	PREFIX_META               = "meta"
	PREFIX_SAMPLE             = "sample"
	MQ_TYPE_ADD_TAG           = "0x00020000"
	MQ_TYPE_DEL_TAG           = "0x00020001"
	MQ_TYPE_DEL_ITEM          = "0x00020002"
	MQ_TYPE_DEL_REP           = "0x00020003"
	STATUS_COOPERATORRING     = "协作中"
	STATUS_COOPERATOR         = "协作"
)

var (
	SUPPLY_STYLE_ALL     = []string{SUPPLY_STYLE_API, SUPPLY_STYLE_BATCH, SUPPLY_STYLE_FLOW}
	SEARCH_DATAITEM_COLS = []string{COL_REPNAME, COL_ITEM_NAME, COL_COMMENT}
)

func getDetetedHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	copy := db.copy()
	defer copy.Close()

	loginName := r.Header.Get("User")
	if loginName == "" {
		return rsp.Json(401, E(ErrorCodeNoLogin))
	}

	dels := []delete{}
	rep := repository{}
	items := []dataItem{}
	var err error
	Q := bson.M{COL_CREATE_USER: loginName}

	repIter := copy.DB(DB_NAME).C(C_REPOSITORY_DEL).Find(Q).Iter()
	defer repIter.Close()
	for repIter.Next(&rep) {
		Q := bson.M{
			COL_CREATE_USER: loginName,
			COL_REPNAME:     rep.Repository_name,
		}
		if items, err = db.getdeletedDataitems(Q); err != nil {
			Log.Infof("get delete dataitem err %s", err.Error())
		}

		delete := delete{
			Rep:   newRepProxy(rep),
			Items: newitemsProxy(items),
		}

		dels = append(dels, delete)
	}

	return rsp.Json(200, E(OK), dels)
}
func createRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, login_name string, l Quota) (int, string) {
	defer db.Close()
	repname := param[PARAM_REP_NAME]
	if err := cheParam(PARAM_REP_NAME, repname); err != nil {
		return rsp.Json(400, err)
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

	if err := cheParam(PARAM_COMMENT_NAME, rep.Comment); err != nil {
		return rsp.Json(400, err)
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
		}{has + 1}
	case ACCESS_PRIVATE:
		max = l.Rep_Private
		opt = struct {
			Private int `json:"private"`
		}{has + 1}
	}

	Log.Infof("createRHandler the result max %d\n", max)
	if has >= max {
		return rsp.Json(400, ErrRepOutOfLimit(max))
	}

	rep.Optime = now.String()
	rep.Ct = now
	rep.Create_user = login_name
	rep.Repository_name = repname
	rep.Items = 0
	rep.CooperateItems = 0
	rep.chkLabel()
	rep.Cooperate = []string{}

	if err := db.DB(DB_NAME).C(C_REPOSITORY).Insert(rep); err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}

	token := r.Header.Get(AUTHORIZATION)
	if _, err := updateUser(login_name, token, opt); err != nil {
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
	relatedItems := strings.TrimSpace(r.FormValue("relatedItems"))
	myRelease := strings.TrimSpace(r.FormValue("myRelease"))

	sortKey, err := getSortKeyByParam("items", showItems)
	if err != nil {
		return rsp.Json(400, ErrInvalidParameter(showItems))
	}

	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}
	rep.Optime = buildTime(rep.Optime)

	if rep.Repaccesstype == ACCESS_PRIVATE {
		Q = bson.M{COL_PERMIT_REPNAME: repname, COL_PERMIT_USER: user}
		if user != "" {
			if rep.Create_user != user && !db.hasPermission(C_REPOSITORY_PERMISSION, Q) && user != "datahub@asiainfo.com" {
				Log.Debugf("[Auth] login name %s, repository name %s.", user, rep.Create_user)
				return rsp.Json(400, E(ErrorCodePermissionDenied))
			}
		} else {
			return rsp.Json(400, E(ErrorCodePermissionDenied))
		}
	}

	var res struct {
		repository
		Cooperate_status string   `json:"cooperatestate,omitempty"`
		Dataitems        []string `json:"dataitems"`
		ItemSize         int      `json:"itemsize"`
	}
	res.repository = rep

	if showItems != "" || relatedItems != "" {

		items := []string{}
		ds := []dataItem{}
		if cooperates, ok := rep.Cooperate.([]interface{}); ok {
			if len(cooperates) > 0 {
				if contains(cooperates, user) {
					res.Cooperate_status = STATUS_COOPERATORRING
				} else {
					res.Cooperate_status = STATUS_COOPERATOR
				}
			}
		}

		Q := bson.M{COL_REPNAME: repname}

		if myRelease != "" && showItems != "" {
			if user != "" && ifCooperate(rep.Cooperate, user) {
				Q[COL_CREATE_USER] = user
			}
		}

		ds, err = db.getDataitems(page_index, page_size, Q, sortKey)
		get(err)

		res.ItemSize = db.countNum(C_DATAITEM, Q)

		for _, v := range ds {
			items = append(items, v.Dataitem_name)
		}
		res.Dataitems = items
	}

	return rsp.Json(200, E(OK), res)
}

func delRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string, msg *Msg) (int, string) {
	defer db.Close()
	repname := strings.TrimSpace(param["repname"])
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}

	if rep.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	if rep.Items > 0 {
		return rsp.Json(400, E(ErrorCodeRepExistDataitem))
	}

	items, err := db.getDataitems(0, SELECT_ALL, Q)
	if err != nil {
		return rsp.Json(400, ErrParseJson(err))
	}

	for _, v := range items {
		if v.Create_user != loginName {
			return rsp.Json(400, ErrRepExistCooperateItem(repname, v.Dataitem_name))
		}
	}

	db.deleteDataitemsFunc(items, msg)

	var opt interface{}
	has := db.countNum(C_REPOSITORY, bson.M{COL_CREATE_USER: r.Header.Get("User"), COL_REP_ACC: rep.Repaccesstype})
	switch rep.Repaccesstype {
	case ACCESS_PUBLIC:
		opt = struct {
			Public int `json:"public"`
		}{has - 1}
	case ACCESS_PRIVATE:
		opt = struct {
			Private int `json:"private"`
		}{has - 1}
	}

	token := r.Header.Get(AUTHORIZATION)
	if _, err := updateUser(loginName, token, opt); err != nil {
		Log.Errorf("create dataitem update User err: %s", err.Error())
	}

	if err := db.delRepository(Q); err != nil {
		return rsp.Json(200, ErrDataBase(err))
	}

	return rsp.Json(200, E(OK))
}

func updateRHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()
	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	Q := bson.M{COL_REPNAME: repname}
	repo, err := db.getRepository(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
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

	u := bson.M{}

	if rep.Comment != "" {
		if err := cheParam(PARAM_COMMENT_NAME, rep.Comment); err != nil {
			return rsp.Json(400, err)
		}
		u[COL_COMMENT] = rep.Comment
	}

	if rep.Repaccesstype != "" {
		if rep.Repaccesstype != ACCESS_PRIVATE && rep.Repaccesstype != ACCESS_PUBLIC {
			return rsp.Json(400, ErrInvalidParameter("repaccesstype"))
		}
		u[COL_REP_ACC] = rep.Repaccesstype
	}

	if u[COL_REP_ACC] != "" {
		token := r.Header.Get(AUTHORIZATION)
		quota := getUserQuota(token, loginName)
		have_pub := db.countNum(C_REPOSITORY, bson.M{COL_CREATE_USER: loginName, COL_REP_ACC: ACCESS_PUBLIC})
		have_pri := db.countNum(C_REPOSITORY, bson.M{COL_CREATE_USER: loginName, COL_REP_ACC: ACCESS_PRIVATE})

		var opt = struct {
			Public  int `json:"public"`
			Private int `json:"private"`
		}{}

		if repo.Repaccesstype == ACCESS_PUBLIC && u[COL_REP_ACC] == ACCESS_PRIVATE {
			if quota.Rep_Private <= have_pri {
				return rsp.Json(400, ErrRepOutOfLimit(quota.Rep_Private))
			}
			users := getSubscribers(Subscripters_By_Rep, repname, "", token)
			if len(users) > 0 {
				for _, user := range users {
					putRepositoryPermission(repname, user, PERMISSION_READ)
				}
			}
			opt.Public = have_pub - 1
			opt.Private = have_pri + 1
		}

		if repo.Repaccesstype == ACCESS_PRIVATE && u[COL_REP_ACC] == ACCESS_PUBLIC {
			if quota.Rep_Public <= have_pub {
				return rsp.Json(400, ErrRepOutOfLimit(quota.Rep_Private))
			}
			exec := bson.M{
				COL_REPNAME:      repname,
				"opt_permission": bson.M{CMD_NOTEQUAL: PERMISSION_WRITE},
			}
			if err := db.delPermit(C_REPOSITORY_PERMISSION, exec); err != nil {
				return rsp.Json(400, ErrDataBase(err))
			}

			opt.Public = have_pub + 1
			opt.Private = have_pri - 1
		}

		if _, err := updateUser(loginName, token, opt); err != nil {
			Log.Errorf("update repository User err: %s", err.Error())
		}
	}

	if len(u) > 0 {
		selector := bson.M{COL_REPNAME: repname}
		u[COL_OPTIME] = time.Now().String()
		update := bson.M{
			CMD_SET: u,
		}
		exec := Execute{
			Collection: C_REPOSITORY,
			Selector:   selector,
			Update:     update,
			Type:       Exec_Type_Update,
		}
		go asynExec(exec)
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
	if (loginName != "" && targetName == "") || (loginName != "" && targetName == loginName) { //login already and search myrepositories
		Q = bson.M{
			CMD_OR: []bson.M{
				bson.M{
					COL_CREATE_USER: loginName, // my repositories
				},
				bson.M{
					COL_REP_COOPERATOR: loginName, // cooperator repositories
				},
			},
		}
	} else if loginName == "" && targetName != "" { // no login nd search targetName
		Q = bson.M{COL_CREATE_USER: targetName, COL_REP_ACC: ACCESS_PUBLIC}
	} else if loginName != "" && targetName != "" && loginName != targetName { //loging and search targetName
		q := []bson.M{}
		p_reps, _ := db.getPermits(C_REPOSITORY_PERMISSION, bson.M{COL_PERMIT_USER: loginName})
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
	}

	rep := []repository{}

	if page_size == -1 {
		if err := db.DB(DB_NAME).C(C_REPOSITORY).Find(Q).Sort("-rank").Select(bson.M{COL_REPNAME: "1", COL_CREATE_USER: "1", COL_REP_COOPERATOR: "1"}).All(&rep); err != nil {
			return rsp.Json(400, ErrDataBase(err))
		}
	} else {
		err := db.DB(DB_NAME).C(C_REPOSITORY).Find(Q).Sort("-rank").Select(bson.M{COL_REPNAME: "1", COL_CREATE_USER: "1", COL_REP_COOPERATOR: "1"}).Skip((page_index - 1) * page_size).Limit(page_size).All(&rep)
		if err != nil {
			return rsp.Json(400, ErrDataBase(err))
		}
	}
	l := []names{}
	for _, v := range rep {
		status := ""
		if cooperates, ok := v.Cooperate.([]interface{}); ok {
			if len(cooperates) > 0 {
				if contains(cooperates, loginName) {
					status = STATUS_COOPERATORRING
				} else {
					status = STATUS_COOPERATOR
				}
			}
		}
		l = append(l, names{Repository_name: v.Repository_name, Cooperate_status: status})
	}

	return rsp.Json(200, E(OK), l)
}

func createDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string) (int, string) {
	defer db.Close()

	repname, itemname := param[PARAM_REP_NAME], param[PARAM_ITEM_NAME]
	if err := cheParam(PARAM_REP_NAME, repname); err != nil {
		return rsp.Json(400, err)
	}
	if err := cheParam(PARAM_ITEM_NAME, itemname); err != nil {
		return rsp.Json(400, err)
	}

	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}

	cooperate := ifCooperate(rep.Cooperate, loginName)
	if rep.Create_user != loginName && !cooperate {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	if n, _ := db.DB(DB_NAME).C(C_DATAITEM).Find(bson.M{COL_REPNAME: repname}).Count(); n > 200 {
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

	if err := cheParam(PARAM_COMMENT_NAME, d.Comment); err != nil {
		return rsp.Json(400, err)
	}

	d.Repository_name = repname
	d.Dataitem_name = itemname
	d.Create_user = loginName
	now := time.Now()
	d.Optime = now.String()
	d.Ct = now
	d.Tags = 0
	d.Cooperate = cooperate

	if d.Itemaccesstype != ACCESS_PRIVATE && d.Itemaccesstype != ACCESS_PUBLIC {
		d.Itemaccesstype = ACCESS_PUBLIC
	}

	if d.Label == "" {
		return rsp.Json(400, ErrNoParameter("label"))
	}

	if err := ifInLabel(d.Label, LABEL_NED_CHECK); err != nil {
		return rsp.Json(400, err)
	}

	if err := chkPrice(d.Price); err != nil {
		return rsp.Json(400, err)
	}
	addPriceElemUid(d.Price)

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

	update := bson.M{"items": 1}
	if rep.Create_user != loginName {
		update["cooperateitems"] = 1
	}

	exec := Execute{
		Collection: C_REPOSITORY,
		Selector:   bson.M{COL_REPNAME: repname},
		Update:     bson.M{CMD_INC: update, CMD_SET: bson.M{COL_OPTIME: now.String()}},
		Type:       Exec_Type_Update,
	}

	go asynExec(exec)

	return rsp.Json(200, E(OK))
}

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
		return rsp.Json(400, ErrDataitemNotFound(itemname))
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
		if err := chkPrice(d.Price); err != nil {
			return rsp.Json(400, err)
		}
		addPriceElemUid(d.Price)
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
		if err := cheParam(PARAM_COMMENT_NAME, d.Comment); err != nil {
			return rsp.Json(400, err)
		}
		u[COL_COMMENT] = d.Comment
	}

	if len(u) > 0 {
		now := time.Now().String()
		u[COL_OPTIME] = now
		update := bson.M{"$set": u}

		exec := []Execute{
			{
				Collection: C_DATAITEM,
				Selector:   selector,
				Update:     update,
				Type:       Exec_Type_Update,
			},
			{
				Collection: C_REPOSITORY,
				Selector:   bson.M{COL_REPNAME: repname},
				Update:     bson.M{"$set": bson.M{COL_OPTIME: now}},
				Type:       Exec_Type_Update,
			},
		}

		go asynExec(exec...)
	}

	if item.Itemaccesstype == ACCESS_PUBLIC && u[COL_ITEM_ACC] == ACCESS_PRIVATE {
		token := r.Header.Get(AUTHORIZATION)
		users := getSubscribers(Subscripters_By_Item, repname, itemname, token)
		if len(users) > 0 {
			for _, user := range users {
				putDataitemPermission(repname, itemname, user)
			}
		}
	}

	if item.Itemaccesstype == ACCESS_PRIVATE && u[COL_ITEM_ACC] == ACCESS_PUBLIC {
		if err := db.delPermit(C_DATAITEM_PERMISSION, selector); err != nil {
			return rsp.Json(400, ErrDataBase(err))
		}
	}

	return rsp.Json(200, E(OK))
}

//curl http://127.0.0.1:8080/repositories/rep123/bear23 -X DELETE -H admin:admin
func delDHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string, msg *Msg) (int, string) {
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
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	if err := db.deleteDataitemFunc(item, msg); err != nil {
		return rsp.Json(200, ErrParseJson(err))
	}

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
		return rsp.Json(400, ErrFieldNotFound("Dataitem Selector", labelname))
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

	update := bson.M{CMD_SET: u}

	exec := Execute{
		Collection: C_SELECT,
		Selector:   selector,
		Update:     update,
		Type:       Exec_Type_Update,
	}

	go asynExec(exec)

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
		return rsp.Json(400, ErrFieldNotFound("Dataitem Selector", labelname))
	}
	return rsp.Json(200, E(OK))
}

func getSelectLabelsHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
	defer db.Close()

	l, ll := []Select{Select{LabelName: "全部精选", Icon: "allselect"}}, []Select{}
	err := db.DB(DB_NAME).C(C_SELECT).Find(nil).Sort("-order").All(&ll)
	if err != nil {
		return rsp.Json(400, ErrDataBase(err))
	}
	l = append(l, ll...)
	return rsp.Json(200, E(OK), l)
}

func createTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string, msg *Msg) (int, string) {
	defer db.Close()
	repname, itemname, tagname := param[PARAM_REP_NAME], param[PARAM_ITEM_NAME], param[PARAM_TAG_NAME]
	if err := cheParam(PARAM_REP_NAME, repname); err != nil {
		return rsp.Json(400, err)
	}
	if err := cheParam(PARAM_ITEM_NAME, itemname); err != nil {
		return rsp.Json(400, err)
	}
	if err := cheParam(PARAM_TAG_NAME, tagname); err != nil {
		return rsp.Json(400, err)
	}

	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q)
	if err == mgo.ErrNotFound {
		return rsp.Json(400, ErrDataitemNotFound(itemname))
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

	exec := []Execute{
		{
			Collection: C_REPOSITORY,
			Selector:   bson.M{COL_REPNAME: repname},
			Update:     bson.M{CMD_SET: bson.M{COL_OPTIME: now}},
			Type:       Exec_Type_Update,
		},
		{
			Collection: C_DATAITEM,
			Selector:   Q,
			Update:     bson.M{CMD_INC: bson.M{"tags": 1}, CMD_SET: bson.M{COL_OPTIME: now}},
			Type:       Exec_Type_Update,
		},
	}

	go asynExec(exec...)

	if msg != nil {
		go func(msg *Msg, t tag) {
			m_t := m_tag{Type: MQ_TYPE_ADD_TAG, Repository_name: t.Repository_name, Dataitem_name: t.Dataitem_name, Tag: t.Tag, Time: t.Optime}
			msg.MqJson(MQ_TOPIC_TO_SUB, m_t)
		}(msg, *t)
	}

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
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	Q_tag := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname, COL_TAG_NAME: tagname}
	if _, err := db.getTag(Q_tag); err == mgo.ErrNotFound {
		return rsp.Json(400, ErrTagNotFound(tagname))
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

	exec := []Execute{
		{
			Collection: C_REPOSITORY,
			Selector:   bson.M{COL_REPNAME: repname},
			Update:     bson.M{CMD_SET: bson.M{COL_OPTIME: now}},
			Type:       Exec_Type_Update,
		},
		{
			Collection: C_DATAITEM,
			Selector:   Q_item,
			Update:     bson.M{CMD_INC: bson.M{"tags": 1}, CMD_SET: bson.M{COL_OPTIME: now}},
			Type:       Exec_Type_Update,
		},
		{
			Collection: C_TAG,
			Selector:   Q_tag,
			Update:     bson.M{"$set": bson.M{COL_COMMENT: t.Comment, COL_OPTIME: now}},
			Type:       Exec_Type_Update,
		},
	}

	go asynExec(exec...)

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
		return rsp.Json(400, ErrTagNotFound(tagname))
	}
	tag.Optime = buildTime(tag.Optime)

	return rsp.Json(200, E(OK), tag)
}

//curl http://127.0.0.1:8080/repositories/NBA/bear23/0001 -H user:admin -X DELETE
func delTagHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB, loginName string, msg *Msg) (int, string) {
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
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	if item.Create_user != loginName {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}

	Q[COL_TAG_NAME] = tagname
	err = db.delTag(Q)
	if err != nil {
		if err == mgo.ErrNotFound {
			if err == mgo.ErrNotFound {
				return rsp.Json(400, ErrTagNotFound(tagname))
			}
		}
		return rsp.Json(400, ErrDataBase(err))
	}

	exec := Execute{
		Collection: C_DATAITEM,
		Selector:   bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname},
		Update:     bson.M{CMD_INC: bson.M{"tags": -1}, CMD_SET: bson.M{COL_OPTIME: time.Now().String()}},
		Type:       Exec_Type_Update,
	}

	go asynExec(exec)

	t := m_tag{Type: MQ_TYPE_DEL_TAG, Repository_name: Q[COL_REPNAME], Dataitem_name: Q[COL_ITEM_NAME], Tag: Q[COL_TAG_NAME], Time: time.Now().String()}
	if msg != nil {
		go func(msg *Msg, t m_tag) {
			msg.MqJson(MQ_TOPIC_TO_SUB, t)
		}(msg, t)
	}
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

	abstract := false
	if p := strings.TrimSpace(r.FormValue("abstract")); p == "1" {
		abstract = true
	}

	Q := bson.M{COL_REPNAME: repname}
	rep, err := db.getRepository(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}

	Q = bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err := db.getDataitem(Q, abstract)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	user := r.Header.Get("User")
	if user != "" {
		if rep.Repaccesstype == ACCESS_PRIVATE {
			if user != rep.Create_user && user != item.Create_user && user != "datahub@asiainfo.com" { //如果是rep的创建者或者是item的创建者,或者是管理员"datahub@asiainfo.com",可以查看私有
				Q := bson.M{COL_PERMIT_REPNAME: rep.Repository_name, COL_PERMIT_USER: user} //如果不是,这需要查看是否在白名单中
				if !db.hasPermission(C_REPOSITORY_PERMISSION, Q) {
					return rsp.Json(400, E(ErrorCodePermissionDenied))
				}
			}
		}
	} else {
		if rep.Repaccesstype == ACCESS_PRIVATE {
			return rsp.Json(400, E(ErrorCodePermissionDenied))
		}
	}

	//if user != "" && ifCooperate(rep.Cooperate, user) && item.Create_user != user {
	//	return rsp.Json(400, E(ErrorCodePermissionDenied)) //虽然是协作者,但是这个item不是自己创建的
	//}

	item.Optime = buildTime(item.Optime)

	var res struct {
		dataItem
		Tags          []tag  `json:"taglist"`
		Permisson     bool   `json:"permission"`
		Stat          string `json:"pricestate"`
		StatCooperate string `json:"cooperatestate"`
	}

	priceStat := getPriceStat(item.Price)
	cooperateStat := getCooperateStat(item, user)

	if abstract == true {
		res.dataItem = item
		res.Stat = priceStat
		res.StatCooperate = cooperateStat
		return rsp.Json(200, E(OK), res)
	}

	if item.Create_user == user {
		res.Permisson = true
	} else {
		switch item.Itemaccesstype {
		case ACCESS_PUBLIC:
			res.Permisson = true
		case ACCESS_PRIVATE:
			Q = bson.M{COL_PERMIT_REPNAME: repname, COL_PERMIT_ITEMNAME: itemname, COL_PERMIT_USER: user}
			res.Permisson = db.hasPermission(C_DATAITEM_PERMISSION, Q)
		}
	}

	b_m, err := db.getFile(PREFIX_META, repname, itemname)
	if err != nil {
		if err != mgo.ErrNotFound {
			Log.Error("get dataitem meta err :", err)
		}
		item.Meta = ""
	} else {
		item.Meta = strings.TrimSpace(string(b_m))
	}

	b_s, err := db.getFile(PREFIX_SAMPLE, repname, itemname)
	if err != nil {
		if err != mgo.ErrNotFound {
			Log.Error("get dataitem sample err :", err)
		}
		item.Sample = ""
	} else {
		item.Sample = strings.TrimSpace(string(b_s))
	}

	Q = bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	tags, err := db.getTags(page_index, page_size, Q)
	get(err)
	buildTagsTime(tags)

	res.dataItem = item
	res.Tags = tags
	res.Stat = priceStat
	res.StatCooperate = cooperateStat

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
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}

	Q = bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	item, err = db.getDataitem(Q)
	if err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrDataitemNotFound(itemname))
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
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	u := bson.M{}
	u["label.sys.select_labels"] = select_labels
	u["label.sys.order"] = order
	update := bson.M{"$set": u}

	exec := Execute{
		Collection: C_DATAITEM,
		Selector:   selector,
		Update:     update,
		Type:       Exec_Type_Update,
	}
	go asynExec(exec)

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
	update := bson.M{CMD_UNSET: u}

	exec := Execute{
		Collection: C_DATAITEM,
		Selector:   selector,
		Update:     update,
		Type:       Exec_Type_Update,
	}

	go asynExec(exec)

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
		private := db.getPermitedReps(username)
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
func searchHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {
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
	Q := bson.M{}
	pub := db.getPublicReps()
	if username != "" {
		private := db.getPermitedReps(username)
		pub = append(pub, private...)
	}

	if len(pub) > 0 {
		Q = bson.M{COL_REPNAME: bson.M{CMD_IN: pub}}
	}

	l := []names{}
	res := map[string]interface{}{}
	text := strings.TrimSpace(r.FormValue("text"))
	if text != "" {
		searchs := strings.Split(text, " ")
		for _, v := range searchs {
			for _, col := range SEARCH_DATAITEM_COLS {
				Query := bson.M{}
				q := bson.M{col: bson.M{CMD_REGEX: v, CMD_OPTION: CMD_CASE_ALL}}
				Query[CMD_AND] = []bson.M{q, Q}
				l := []search{}
				db.DB(DB_NAME).C(C_DATAITEM).Find(Query).Sort("-optime").Select(bson.M{COL_REPNAME: "1", COL_ITEM_NAME: "1", "optime": "1"}).All(&l)
				for _, v := range l {
					if sc, ok := res[fmt.Sprintf("%s/%s", v.Repository_name, v.Dataitem_name)]; ok {
						sc.(*score).matchCount++
					} else {
						sc := score{optime: fmt.Sprintf("%s", v.Optime), matchCount: 1}
						res[fmt.Sprintf("%s/%s", v.Repository_name, v.Dataitem_name)] = &sc
					}
				}
			}
		}

		res_reverse, res_reverse_2, res_reverse_3 := make(Ms), make(Ms), make(Ms)

		for k, v := range res {
			sc := v.(*score)
			switch sc.matchCount {
			case 1:
				res_reverse[v.(*score).optime] = k
			case 2:
				res_reverse_2[v.(*score).optime] = k
			case 3:
				res_reverse_3[v.(*score).optime] = k
			}
		}

		res_reverse_3.sortMapToArray(&l)
		res_reverse_2.sortMapToArray(&l)
		res_reverse.sortMapToArray(&l)

	} else {
		Q := bson.M{COL_REPNAME: bson.M{CMD_IN: pub}}
		db.DB(DB_NAME).C(C_DATAITEM).Find(Q).Sort("-rank").Select(bson.M{COL_REPNAME: "1", COL_ITEM_NAME: "1", "ct": "1"}).All(&l)
	}

	length := len(l)
	result := struct {
		Results []names `json:"results"`
		Total   int     `json:"total"`
	}{
		l,
		length,
	}

	if length < page_index*page_size && length >= (page_index-1)*page_size {
		result.Results = l[(page_index-1)*page_size : length]
	} else if length < page_index*page_size {
		result.Results = l
	} else if length >= page_index*page_size {
		result.Results = l[(page_index-1)*page_size : page_index*page_size]
	}

	result.Total = length

	return rsp.Json(200, E(OK), result)
}

func (m Ms) sortMapToArray(l *[]names) {
	keys := StringSlice{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		str := strings.Split(m[k].(string), "/")
		*l = append(*l, names{Repository_name: str[0], Dataitem_name: str[1]})
	}
}

type score struct {
	optime     string
	matchCount int
}
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] >= p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
