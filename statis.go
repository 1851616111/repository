package main

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

func getStatisHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	items, err := db.getDataitems(0, ALL_DATAITEMS, nil)
	get(err)

	log.Printf("statis datitem total %d", len(items))
	for _, v := range items {
		Q := bson.M{COL_REPNAME: v.Repository_name, COL_ITEM_NAME: v.Dataitem_name}
		n, err := db.DB(DB_NAME).C(C_TAG).Find(Q).Count()
		get(err)
		if n < v.Tags {
			log.Println("correct %s/%s tags = %d", v.Repository_name, v.Dataitem_name, n)
		}
		exec := bson.M{CMD_SET: bson.M{COL_ITEM_TAGS: n}}
		go asynUpdateOpt(C_DATAITEM, Q, exec)
	}
	log.Println("statis datitem over")
	return rsp.Json(200, E(OK))
}
