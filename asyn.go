package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
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
		db.Clone()
		copy := DB{*db.Copy()}
		go copy.handle(exec)
	}
}

func (db *DB) handle(e exec) {
	err := db.DB(DB_NAME).C(e.collectionName).Update(e.selector, e.update)
	if err != nil {
		log.Println("queue handle execute update err ", err)
	}
	db.Close()
}

func asynOpt(collection string, selector, updater bson.M) {
	q_c.producer(exec{collection, selector, updater})
}
