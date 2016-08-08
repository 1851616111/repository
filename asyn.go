package main

import (
	"gopkg.in/mgo.v2"
	"time"
)
// 这是一个批量更新的异步处理，负责与mongodb批量更新。
// 负责更新（update）,创建（insert）,不存在时创建，存在时更新（upsert）

const (
	CHANNE_MAX_ZISE      = 20000
	Exec_Type_Update     = "update"
	Exec_Type_Upsert     = "upsert"
	Exec_Type_Insert     = "insert"
	Asyn_Interval_Update = 500
)

var (
	queue = make(chan Execute, CHANNE_MAX_ZISE)
)

type Execute struct {
	Collection string
	Selector   interface{}
	Update     interface{}
	Type       string
}

type execs []Execute

type Queue struct {
	q chan Execute
}

//向需要异步更新的操作提供输入方法
func (q *Queue) producer(e ...Execute) {
	for _, v := range e {
		queue <- v
	}
}

//对于insert和upsert操作，立即进行mongodb操作
//对于update操作，每隔500ms进行更新
func (q *Queue) serve(db *DB) {
	updates, upserts, inserts := execs{}, execs{}, execs{}
	for {
		select {
		case exec := <-q.q:
			switch exec.Type {
			case Exec_Type_Update:
				updates = append(updates, exec)
			case Exec_Type_Upsert:
				upserts = append(upserts, exec)
			case Exec_Type_Insert:
				inserts = append(inserts, exec)
			}

			if len(inserts) > 0 {
				inserts.serve(db, Exec_Type_Insert)
			}

			if len(upserts) > 0 {
				upserts.serve(db, Exec_Type_Upsert)
			}

		case <-time.After(time.Millisecond * Asyn_Interval_Update):
			updates.serve(db, Exec_Type_Update)
		}
	}
}

func (db *DB) bulkHandle(es *execs, Type string) {
	defer db.Close()
	m := make(M)
	for _, exec := range *es {
		bulk, ok := m[exec.Collection]
		if !ok {
			bulk = db.DB(DB_NAME).C(exec.Collection).Bulk()
			m[exec.Collection] = bulk
		}

		switch Type {
		case Exec_Type_Update:
			m[exec.Collection].(*mgo.Bulk).Update(exec.Selector, exec.Update)
		case Exec_Type_Upsert:
			m[exec.Collection].(*mgo.Bulk).Upsert(exec.Selector, exec.Update)
		case Exec_Type_Insert:
			m[exec.Collection].(*mgo.Bulk).Insert(exec.Update)
		}
	}

	for _, bulk := range m {
		b := bulk.(*mgo.Bulk)
		b.Unordered()
		res, err := b.Run()
		if err != nil {
			Log.Errorf("queue asyn operator err %s. result %+v", err.Error(), res)
		}
		//todo save fail bulk and retry
	}

	*es = execs{}
}

func (p *execs) serve(db *DB, Type string) {
	if len(*p) > 0 {
		copy := db.copy()
		go copy.bulkHandle(p, Type)
	}
}

func asynExec(e ...Execute) {
	q_c.producer(e...)
}
