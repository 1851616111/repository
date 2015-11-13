package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
)

const (
	MONGO_REGEX = "$regex"
)

var SEARCH_DATAITEM_COLS = []string{COL_REP_NAME, COL_ITEM_NAME}

//http://127.0.0.1:8080/search -d "text=123 123 14"
func searchHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {

	page_index, page_size := PAGE_INDEX, PAGE_SIZE
	if p := strings.TrimSpace(r.FormValue("page")); p != "" {
		if page_index, _ = strconv.Atoi(p); page_index <= 0 {
			return rsp.Json(400, ErrInvalidParameter("page"))
		}

	}
	if p := strings.TrimSpace(r.FormValue("size")); p != "" {
		if page_size, _ = strconv.Atoi(p); page_size <= 0 {
			return rsp.Json(400, ErrInvalidParameter("size"))
		}
	}
	res := map[string]interface{}{}
	text := strings.TrimSpace(r.PostFormValue("text"))
	searchs := strings.Split(text, " ")
	for _, v := range searchs {
		for _, col := range SEARCH_DATAITEM_COLS {
			l := []names{}
			db.DB(DB_NAME).C(C_DATAITEM).Find(bson.M{col: bson.M{"$regex": v}}).Select(bson.M{COL_REP_NAME: "1", COL_ITEM_NAME: "1", "ct": "1"}).All(&l)
			for _, v := range l {
				res[fmt.Sprintf("%s/%s", v.Repository_name, v.Dataitem_name)] = 1
			}
		}
	}

	l := []names{}
	for k, _ := range res {
		str := strings.Split(k, "/")
		l = append(l, names{str[0], str[1]})
	}
	return rsp.Json(200, E(OK), l[(PAGE_INDEX-1)*PAGE_SIZE:PAGE_INDEX*PAGE_SIZE], len(l))
}