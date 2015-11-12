package main

import "gopkg.in/mgo.v2/bson"

func (db *DB) getRepository(query bson.M) (repository, error) {
	res := new(repository)
	err := db.DB(DB_NAME).C(C_REPOSITORY).Find(query).One(res)
	return *res, err
}

func (db *DB) delRepository(exec bson.M)  error {
	return db.DB(DB_NAME).C(C_REPOSITORY).Remove(exec)
}
