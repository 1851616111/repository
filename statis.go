package main

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func getDStatisHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	items, err := db.getDataitems(0, ALL_DATAITEMS, nil)
	get(err)

	log.Infof("statis datitem total %d", len(items))
	for _, v := range items {
		Q := bson.M{COL_REPNAME: v.Repository_name, COL_ITEM_NAME: v.Dataitem_name}
		n, err := db.DB(DB_NAME).C(C_TAG).Find(Q).Count()
		get(err)
		if n < v.Tags {
			log.Infof("correct %s/%s tags = %d", v.Repository_name, v.Dataitem_name, n)
			exec := bson.M{CMD_SET: bson.M{COL_ITEM_TAGS: n}}
			go asynUpdateOpt(C_DATAITEM, Q, exec)
		}
	}
	log.Info("statis datitem over")
	return rsp.Json(200, E(OK))
}

func getRStatisHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	reps, err := db.getRepositories(nil)
	get(err)

	log.Infof("statis repository total %d", len(reps))
	for _, v := range reps {
		Q := bson.M{COL_REPNAME: v.Repository_name}
		n, err := db.DB(DB_NAME).C(C_DATAITEM).Find(Q).Count()
		get(err)
		if n < v.Items {
			log.Infof("correct %s items = %d", v.Repository_name, n)
			exec := bson.M{CMD_SET: bson.M{COL_REP_ITEMS: n}}
			go asynUpdateOpt(C_REPOSITORY, Q, exec)
		}
	}
	log.Error("statis repository over")
	return rsp.Json(200, E(OK))
}

func getStatisHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	getDStatisHandler(r, rsp, param, db)
	getRStatisHandler(r, rsp, param, db)
	return rsp.Json(200, E(OK))
}
