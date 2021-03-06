package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//查询dataitem的所有tags数量，更新到dataitem的tag数量字段中
func (db *DB) getDStatis() {
	items, err := db.getDataitems(0, SELECT_ALL, nil)
	get(err)
	execs := []Execute{}

	Log.Infof("statis datitem total %d", len(items))
	for _, v := range items {
		Q := bson.M{COL_REPNAME: v.Repository_name, COL_ITEM_NAME: v.Dataitem_name}
		n, err := db.DB(DB_NAME).C(C_TAG).Find(Q).Count()
		get(err)
		if n != v.Tags {
			Log.Infof("correct %s/%s tags = %d", v.Repository_name, v.Dataitem_name, n)
			exec := Execute{
				Collection: C_DATAITEM,
				Selector:   Q,
				Update:     bson.M{CMD_SET: bson.M{COL_ITEM_TAGS: n}},
				Type:       Exec_Type_Update,
			}
			execs = append(execs, exec)
		}
	}
	go asynExec(execs...)
	Log.Info("statis datitem over")
}

//查询repository的所有dataitem数量，更新到repository的dataitem数量字段中
func (db *DB) getRStatis() {
	reps, err := db.getRepositories(nil)
	get(err)
	execs := []Execute{}

	Log.Infof("statis repository total %d", len(reps))
	for _, v := range reps {
		Q := bson.M{COL_REPNAME: v.Repository_name}
		n, err := db.DB(DB_NAME).C(C_DATAITEM).Find(Q).Count()
		get(err)
		if n != v.Items {
			Log.Infof("correct %s items = %d", v.Repository_name, n)

			exec := Execute{
				Collection: C_REPOSITORY,
				Selector:   Q,
				Update:     bson.M{CMD_SET: bson.M{COL_REP_ITEMS: n}},
				Type:       Exec_Type_Update,
			}
			execs = append(execs, exec)
		}
	}

	go asynExec(execs...)
	Log.Info("statis repository over")
}

func staticLoop(db *DB) {

	for {
		copy := DB{*db.Copy()}
		defer copy.Close()
		time.Sleep(time.Hour)
		copy.getDStatis()
		copy.getRStatis()
	}
}

type statis struct {
	UserNum int `json:"users"`
	RepNum  int `json:"reps"`
	ItemNum int `json:"items"`
	TagNum  int `json:"tags"`
}
