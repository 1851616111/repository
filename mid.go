package main

import "gopkg.in/mgo.v2/bson"

func (db *DB) getRepository(query bson.M) (repository, error) {
	res := new(repository)
	err := db.DB(DB_NAME).C(C_REPOSITORY).Find(query).One(res)
	return *res, err
}

func (db *DB) getRepositories(query bson.M) ([]repository, error) {
	res := []repository{}
	err := db.DB(DB_NAME).C(C_REPOSITORY).Find(query).All(res)
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
	err := db.DB(DB_NAME).C(C_DATAITEM).Find(query).All(res)
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

func (db *DB) getPermitByUser(userName string) ([]Repository_Permit, error) {
	l := []Repository_Permit{}
	Q := bson.M{"user_name": userName}
	if err := db.DB(DB_NAME).C(C_REPOSITORY_PERMIT).Find(Q).All(&l); err != nil {
		return nil, err
	}
	return l, nil
}
