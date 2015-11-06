package main

import (
	"errors"
	"github.com/go-martini/martini"
	"net/http"
)

//curl http://127.0.0.1:8080/repositories/NBA/bear
func getDMid(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (*DataItem, error) {
	repname := param["repname"]
	if repname == "" {
		return nil, errors.New("no param repname")
	}
	itemname := param["itemname"]
	if itemname == "" {
		return nil, errors.New("no param itemname")
	}
	d := &DataItem{Repository_name: repname, Dataitem_name: itemname}
	exists, err := db.getDataitem(d)
	if exists {
		return d, nil
	} else {
		return nil, err
	}
}
