package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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
	l := []exec{}
	for {
		select {
		case exec := <-q.q:
			l = append(l, exec)
		case <-time.After(time.Second):
			if len(l) > 0 {
				copy := db.copy()
				go copy.bulkHandle(&l)
			}
		}
	}
}

func (db *DB) bulkHandle(es *[]exec) {
	defer db.Close()
	m := make(M)
	for _, exec := range *es {
		bulk, ok := m[exec.collectionName]
		if !ok {
			b := db.DB(DB_NAME).C(exec.collectionName).Bulk()
			m[exec.collectionName] = b
			b.Update(exec.selector, exec.update)
		} else {
			bulk.(*mgo.Bulk).Update(exec.selector, exec.update)
		}
	}

	for _, bulk := range m {
		b := bulk.(*mgo.Bulk)
		b.Unordered()
		res, err := b.Run()
		if err != nil {
			Log.Errorf("queue asyn operator err %s. result %+v", err.Error(), res)
		}
	}

	*es = []exec{}
}

func asynUpdateOpt(collection string, selector, updater bson.M) {
	q_c.producer(exec{collection, selector, updater})
}
