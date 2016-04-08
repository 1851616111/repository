package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	sorts = map[string]map[string]string{
		"items": map[string]string{
			"update_time_up":   "optime",
			"update_time_down": "-optime",
			"rank_up":          "rank",
			"rank_down":        "-rank",
			"1":                "-rank",
			"":                 "-rank",
		}}
)

const (
	C_FS          = "datahub_fs"
	MAX_FILE_SIZE = 800000
	SELECT_ALL    = -1
)

func (db *DB) getRepository(query bson.M) (repository, error) {
	res := new(repository)
	err := db.DB(DB_NAME).C(C_REPOSITORY).Find(query).One(res)
	return *res, err
}

func (db *DB) getRepositories(query bson.M) ([]repository, error) {
	res := []repository{}
	err := db.DB(DB_NAME).C(C_REPOSITORY).Find(query).All(&res)
	return res, err
}

func (db *DB) delRepository(exec bson.M) error {
	result := repository{}
	if err := db.DB(DB_NAME).C(C_REPOSITORY).Find(exec).One(&result); err != nil {
		Log.Error("del repository %#v err %s", exec, err.Error())
		return err
	}
	result.St = time.Now()

	e := Execute{
		Collection: C_REPOSITORY_DEL,
		Update:     result,
		Type:       Exec_Type_Insert,
	}
	go asynExec(e)

	return db.DB(DB_NAME).C(C_REPOSITORY).Remove(exec)
}

func (db *DB) getDataitem(query bson.M, abstract ...bool) (dataItem, error) {
	res := new(dataItem)
	var err error

	if len(abstract) > 0 && abstract[0] == true {
		err = db.DB(DB_NAME).C(C_DATAITEM).Find(query).Select(bson.M{COL_COMMENT: "1", COL_CREATE_USER: "1", COL_ITEM_TAGS: "1", COL_LABEL: "1", COL_ITEM_COOPERATOR: "1"}).One(res)
	}
	err = db.DB(DB_NAME).C(C_DATAITEM).Find(query).One(res)

	return *res, err
}

func (db *DB) getDataitems(pageIndex, pageSize int, query bson.M, sortBy ...string) ([]dataItem, error) {
	sort := "-rank"
	if len(sortBy) > 0 {
		sort = sortBy[0]
	}

	res := []dataItem{}
	var err error
	if pageSize == SELECT_ALL {
		err = db.DB(DB_NAME).C(C_DATAITEM).Find(query).Sort(sort).All(&res)
	} else {
		err = db.DB(DB_NAME).C(C_DATAITEM).Find(query).Sort(sort).Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&res)
	}
	return res, err
}

func getSortKeyByParam(paramName, paramValue string) (string, error) {
	paramSort, ok := sorts[paramName]
	if !ok {
		return "", errors.New("no found param " + paramName)
	}
	sortString, ok := paramSort[paramValue]
	if !ok {
		return "", errors.New("wrong param value")
	}

	return sortString, nil
}

func (db *DB) getdeletedDataitems(query bson.M) ([]dataItem, error) {
	res := []dataItem{}
	err := db.DB(DB_NAME).C(C_DATAITEM_DEL).Find(query).Sort("-ct").All(&res)
	return res, err
}

func (db *DB) delDataitem(exec bson.M) error {
	result := dataItem{}
	if err := db.DB(DB_NAME).C(C_DATAITEM).Find(exec).One(&result); err != nil {
		Log.Error("del dataitem %#v err %s", exec, err.Error())
		return err
	}

	e := Execute{
		Collection: C_DATAITEM_DEL,
		Update:     result,
		Type:       Exec_Type_Insert,
	}
	go asynExec(e)

	return db.DB(DB_NAME).C(C_DATAITEM).Remove(exec)
}

func (db *DB) countNum(collection string, query bson.M) (i int) {
	i, _ = db.DB(DB_NAME).C(collection).Find(query).Count()
	return
}

func (db *DB) countUsers(collection string, query bson.M) (i int) {
	users := []string{}
	db.DB(DB_NAME).C(collection).Find(query).Distinct(COL_CREATE_USER, &users)
	return len(users)
}

func (db *DB) getTag(query bson.M) (tag, error) {
	res := new(tag)
	err := db.DB(DB_NAME).C(C_TAG).Find(query).One(&res)
	return *res, err
}

func (db *DB) getTags(pageIndex, pageSize int, query bson.M) ([]tag, error) {
	res := []tag{}
	var err error

	if pageSize == SELECT_ALL {
		err = db.DB(DB_NAME).C(C_TAG).Find(query).All(&res)
	} else {
		err = db.DB(DB_NAME).C(C_TAG).Find(query).Sort("-optime").Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&res)
	}
	return res, err
}

func (db *DB) delTag(exec bson.M) error {
	err := db.DB(DB_NAME).C(C_TAG).Remove(exec)
	return err
}

func (db *DB) delSelect(exec bson.M) error {
	err := db.DB(DB_NAME).C(C_SELECT).Remove(exec)
	return err
}

func (db *DB) getSelect(query bson.M) (Select, error) {
	res := new(Select)
	err := db.DB(DB_NAME).C(C_SELECT).Find(query).One(&res)
	return *res, err
}

func (db *DB) getPermits(collection string, query bson.M, page ...[]int) (interface{}, error) {
	var err error
	var list interface{}

	if len(page) > 0 && len(page[0]) == 2 && page[0][1] != -1 {
		pageIndex := page[0][0]
		pageSize := page[0][1]
		switch collection {
		case C_DATAITEM_PERMISSION:
			l := []Item_Permission{}
			err = db.DB(DB_NAME).C(collection).Find(query).Sort(COL_PERMIT_USER).Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&l)
			list = l
		case C_REPOSITORY_PERMISSION:
			l := []Rep_Permission{}
			err = db.DB(DB_NAME).C(collection).Find(query).Sort(COL_PERMIT_USER).Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&l)
			list = l
		}

	} else {
		switch collection {
		case C_DATAITEM_PERMISSION:
			l := []Item_Permission{}
			err = db.DB(DB_NAME).C(collection).Find(query).All(&l)
			list = l
		case C_REPOSITORY_PERMISSION:
			l := []Rep_Permission{}
			err = db.DB(DB_NAME).C(collection).Find(query).All(&l)
			list = l
		}
	}

	if err != nil {
		return list, err
	}

	return list, nil
}

func (db *DB) countPermits(collection string, query bson.M) (int, error) {
	return db.DB(DB_NAME).C(collection).Find(query).Count()
}

func (db *DB) delPermit(collection string, exec bson.M) (err error) {
	_, err = db.DB(DB_NAME).C(collection).RemoveAll(exec)
	return
}

func (db *DB) countCooperator(repName string) (int, error) {
	return db.countPermits(C_REPOSITORY_PERMISSION, bson.M{COL_REPNAME: repName, "opt_permission": PERMISSION_WRITE})
}

func (db *DB) delRepCooperator(repName string) error {
	err := db.delPermit(C_REPOSITORY_PERMISSION, bson.M{COL_REPNAME: repName, "opt_permission": PERMISSION_WRITE})
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	return nil
}

func setFileName(prefix, repname, itemname string) string {
	return fmt.Sprintf("%s_%s_%s", prefix, repname, itemname)
}

func (db *DB) setFile(prefix, repName, itemName string, b []byte) *Error {
	f, err := db.DB(DB_NAME).GridFS(C_FS).Create("")
	if err != nil {
		return ErrFile(err)
	}
	_, err = f.Write(b)
	if err != nil {
		return ErrFile(err)
	}
	f.SetMeta(bson.M{"prefix": prefix, COL_REPNAME: repName, COL_ITEM_NAME: itemName})
	f.SetName(setFileName(prefix, repName, itemName))
	err = f.Close()
	if err != nil {
		return ErrFile(err)
	}
	return nil
}

func (db *DB) getFile(prefix, repName, itemName string) ([]byte, error) {
	b := make([]byte, MAX_FILE_SIZE)
	file, err := db.DB(DB_NAME).GridFS(C_FS).Open(setFileName(prefix, repName, itemName))
	if err != nil {
		return b, err
	}
	n, err := file.Read(b)
	if err != nil {
		return b, err
	}
	err = file.Close()
	if err != nil {
		return b, err
	}
	return b[:n], nil
}

func (db *DB) delFile(prefix, repName, itemName string) *Error {
	file, err := db.getFile(prefix, repName, itemName)
	if err != nil {
		return ErrFile(err)
	}

	newPrefix := fmt.Sprintf("del_%s", prefix)
	if err := db.setFile(newPrefix, repName, itemName, file); err != nil {
		return err
	}

	if err := db.DB(DB_NAME).GridFS(C_FS).Remove(setFileName(prefix, repName, itemName)); err != nil {
		return ErrFile(err)
	}

	return nil
}

func (db *DB) hasPermission(collection string, query bson.M) bool {
	n, _ := db.DB(DB_NAME).C(collection).Find(query).Count()
	switch n {
	case 0:
		return false
	case 1:
		return true
	default:
		Log.Errorf("query %s  total=%n invalid", collection, n)
		return true
	}
}

func buildTagsTime(tags []tag) {
	for i, v := range tags {
		tags[i].Optime = buildTime(v.Optime)
	}
}

func (db *DB) getPermitedReps(userName string) []string {
	l := []string{}
	if userName != "" {
		p_reps, _ := db.getPermits(C_REPOSITORY_PERMISSION, bson.M{COL_PERMIT_USER: userName})
		if l_p_reps, ok := p_reps.([]Rep_Permission); ok {
			if len(l_p_reps) > 0 {
				for _, v := range p_reps.([]Rep_Permission) {
					l = append(l, v.Repository_name)
				}
			}
		}

		my_pri_reps, _ := db.getRepositories(bson.M{COL_REP_ACC: ACCESS_PRIVATE, COL_CREATE_USER: userName})
		for _, v := range my_pri_reps {
			l = append(l, v.Repository_name)
		}

	}
	return l
}

func (db *DB) getPublicReps() []string {
	s := []string{}
	if l, _ := db.getRepositories(bson.M{COL_REP_ACC: ACCESS_PUBLIC}); len(l) > 0 {
		for _, v := range l {
			s = append(s, v.Repository_name)
		}
	}
	return s
}

func (db *DB) deleteDataitemsFunc(items []dataItem, msg *Msg) {
	for i, _ := range items {
		if err := db.deleteDataitemFunc(items[i], msg); err != nil {
			Log.Errorf("delete dataitems %#v  err %s", items[i], err)
		}
	}
}

func (db *DB) deleteDataitemFunc(item dataItem, msg *Msg) error {
	Q := bson.M{COL_REPNAME: item.Repository_name, COL_ITEM_NAME: item.Dataitem_name}

	if err := db.delDataitem(Q); err != nil {
		return err
	}

	go func(db *DB, repname, itemname string) {
		defer db.Close()
		db.delFile(PREFIX_META, repname, itemname)
		db.delFile(PREFIX_SAMPLE, repname, itemname)
	}(db.copy(), item.Repository_name, item.Dataitem_name)

	exec := Execute{
		Collection: C_REPOSITORY,
		Selector:   bson.M{COL_REPNAME: item.Repository_name},
		Update:     bson.M{CMD_INC: bson.M{"items": -1, "cooperateitems": -1}, CMD_SET: bson.M{COL_OPTIME: time.Now().String()}},
		Type:       Exec_Type_Update,
	}

	go asynExec(exec)

	tmp := m_item{TimeId: fmt.Sprintf("%d", item.Ct.Unix()), Type: MQ_TYPE_DEL_ITEM, Repository_name: Q[COL_REPNAME], Dataitem_name: Q[COL_ITEM_NAME], Time: time.Now().String()}
	if msg != nil {
		go func(msg *Msg, item m_item) {
			msg.MqJson(MQ_TOPIC_TO_SUB, tmp)
		}(msg, tmp)
	}

	return nil
}
