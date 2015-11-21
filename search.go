package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	CMD_REGEX    = "$regex"
	CMD_OPTION   = "$options"
	CMD_AND      = "$and"
	CMD_CASE_ALL = "$i"
)

var SEARCH_DATAITEM_COLS = []string{COL_REPNAME, COL_ITEM_NAME, COL_COMMENT}

//curl http://127.0.0.1:8089/search?text="123 123 14"
func searchHandler(r *http.Request, rsp *Rsp, db *DB) (int, string) {

	page_index, page_size := PAGE_INDEX, PAGE_SIZE_SEARCH
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

	username := r.Header.Get("User")
	Q := bson.M{}
	pub := db.getPublicReps()
	if username != "" {
		private := db.getPrivateReps(username)
		pub = append(pub, private...)
	}

	if len(pub) > 0 {
		Q = bson.M{COL_REPNAME: bson.M{CMD_IN: pub}}
	}

	l := []names{}
	res := map[string]interface{}{}
	text := strings.TrimSpace(r.FormValue("text"))
	if text != "" {
		searchs := strings.Split(text, " ")
		for _, v := range searchs {
			for _, col := range SEARCH_DATAITEM_COLS {
				Query := bson.M{}
				q := bson.M{col: bson.M{CMD_REGEX: v, CMD_OPTION: CMD_CASE_ALL}}
				Query[CMD_AND] = []bson.M{q, Q}
				l := []search{}
				db.DB(DB_NAME).C(C_DATAITEM).Find(Query).Sort("-ct").Select(bson.M{COL_REPNAME: "1", COL_ITEM_NAME: "1", "ct": "1"}).All(&l)
				for _, v := range l {
					res[fmt.Sprintf("%s/%s", v.Repository_name, v.Dataitem_name)] = fmt.Sprintf("%d", v.Ct.Unix())
				}
			}
		}
		res_reverse := map[string]interface{}{}
		for k, v := range res {
			res_reverse[v.(string)] = k
		}

		var keys []string
		for k := range res_reverse {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			str := strings.Split(res_reverse[k].(string), "/")
			l = append(l, names{str[0], str[1]})
		}
	} else {
		Q = bson.M{COL_REPNAME: bson.M{CMD_IN: l}}
		db.DB(DB_NAME).C(C_DATAITEM).Find(Q).Limit(PAGE_SIZE_SEARCH).Sort("-ct").Select(bson.M{COL_REPNAME: "1", COL_ITEM_NAME: "1", "ct": "1"}).All(&l)
	}

	length := len(l)
	result := struct {
		Results []names `json:"results"`
		Total   int     `json:"total"`
	}{
		l,
		length,
	}

	if length < page_index*page_size && length >= (page_index-1)*page_size {
		result.Results = l[(page_index-1)*page_size : length]
	} else if length < page_index*page_size {
		result.Results = l
	} else if length >= page_index*page_size {
		result.Results = l[(page_index-1)*page_size : page_index*page_size]
	}

	if length < page_index*page_size && length >= (page_index-1)*page_size {
		result.Results = l[(page_index-1)*page_size : length]
	} else if length < page_index*page_size {
		result.Results = l
	} else if length >= page_index*page_size {
		result.Results = l[(page_index-1)*page_size : page_index*page_size]
	}

	result.Total = length

	return rsp.Json(200, E(OK), result)
}
