package main

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	CHANNE_MAX_ZISE = 20000
)

var (
	queueChannel = make(chan exec, CHANNE_MAX_ZISE)
)

type exec struct {
	collectionName string
	selector       bson.M
	update         bson.M
}

type Queue struct {
	q chan exec
}

func (q *Queue) producer(e exec) {
	queueChannel <- e
}

func (q *Queue) serve(db *DB) {
	for {
		exec := <-q.q
		copy := db.copy()
		go copy.handle(exec)
	}
}

func (db *DB) handle(e exec) {
	err := db.DB(DB_NAME).C(e.collectionName).Update(e.selector, e.update)
	if err != nil {
		Log.Errorf("queue asyn opt err %s. select :: %+v execute :: %+v", err.Error(), e.selector, e.update)
	}
	db.Close()
}

func asynUpdateOpt(collection string, selector, updater bson.M) {
	q_c.producer(exec{collection, selector, updater})
}
