package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	C_FS          = "datahub_fs"
	MAX_FILE_SIZE = 8192
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
	return db.DB(DB_NAME).C(C_REPOSITORY).Remove(exec)
}

func (db *DB) getDataitem(query bson.M) (dataItem, error) {
	res := new(dataItem)
	err := db.DB(DB_NAME).C(C_DATAITEM).Find(query).One(res)
	return *res, err
}

func (db *DB) getDataitems(pageIndex, pageSize int, query bson.M) ([]dataItem, error) {
	res := []dataItem{}
	var err error
	if pageSize == -1 {
		err = db.DB(DB_NAME).C(C_DATAITEM).Find(query).Sort("-ct").Select(bson.M{COL_ITEM_NAME: "1"}).All(&res)
	} else {
		err = db.DB(DB_NAME).C(C_DATAITEM).Find(query).Sort("-ct").Select(bson.M{COL_ITEM_NAME: "1"}).Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&res)
	}
	return res, err
}

func (db *DB) delDataitem(exec bson.M) error {
	return db.DB(DB_NAME).C(C_DATAITEM).Remove(exec)
}

func (db *DB) getTag(query bson.M) (tag, error) {
	res := new(tag)
	err := db.DB(DB_NAME).C(C_TAG).Find(query).One(&res)
	return *res, err
}

func (db *DB) getTags(pageIndex, pageSize int, query bson.M) ([]tag, error) {
	res := []tag{}
	err := db.DB(DB_NAME).C(C_TAG).Find(query).Sort("-optime").Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&res)
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

func (db *DB) getPermits(collection string, query bson.M) (interface{}, error) {
	var err error
	switch collection {
	case C_DATAITEM_PERMISSION:
		l := []Item_Permission{}
		err = db.DB(DB_NAME).C(collection).Find(query).All(&l)
		if err != nil {
			return l, err
		}
		log.Println(query)
		log.Println(l)
		return l, nil
	case C_REPOSITORY_PERMISSION:
		l := []Rep_Permission{}
		err = db.DB(DB_NAME).C(collection).Find(query).All(&l)
		if err != nil {
			return l, err
		}
		log.Println(query)
		log.Println(l)
		return l, nil
	}
	return nil, errors.New("unknow err")
}

func (db *DB) delPermit(collection string, exec bson.M) error {
	log.Println(exec)
	return db.DB(DB_NAME).C(collection).Remove(exec)
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
		log.Printf("query %s  total=%n invalid", collection, n)
		return true
	}
}

func buildTagsTime(tags []tag) {
	for i, v := range tags {
		tags[i].Optime = buildTime(v.Optime)
	}
}
