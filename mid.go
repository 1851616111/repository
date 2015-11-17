package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const (
	C_FS = "datahub_fs"
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

func (db *DB) getDataitems(query bson.M) ([]dataItem, error) {
	res := []dataItem{}
	err := db.DB(DB_NAME).C(C_DATAITEM).Find(query).All(&res)
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

func (db *DB) getTags(query bson.M) ([]tag, error) {
	res := []tag{}
	err := db.DB(DB_NAME).C(C_TAG).Find(query).All(&res)
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

func (db *DB) getPermitByUser(query bson.M) ([]Repository_Permit, error) {
	l := []Repository_Permit{}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMIT).Find(query).All(&l); err != nil {
		return nil, err
	}
	return l, nil
}

func (db *DB) setFile(repname, itemname string, b []byte) (string, *Error) {
	f, err := db.DB(DB_NAME).GridFS(C_FS).Create("")
	if err != nil {
		return "", ErrFile(err)
	}
	_, err = f.Write(b)
	if err != nil {
		return "", ErrFile(err)
	}
	f.SetMeta(bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname})
	f.SetName(fmt.Sprintf("%s/%s", repname, itemname))
	err = f.Close()
	if err != nil {
		return "", ErrFile(err)
	}
	return f.Id().(bson.ObjectId).Hex(), nil
}

func (db *DB) getFile(fileName string) ([]byte, error) {
	file, err := db.DB(DB_NAME).GridFS(C_FS).Open(fileName)
	get(err)
	b := make([]byte, 8192)
	_, err = file.Read(b)
	get(err)
	fmt.Println(string(b))
	get(err)
	err = file.Close()
	get(err)

	return b, nil
}
