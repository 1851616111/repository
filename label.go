package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"regexp"
)

const (
	LABEL_OPT_ADD = CMD_SET
	LABEL_OPT_SUB = CMD_UNSET
)

var (
	LABEL_OWNER *regexp.Regexp
	LABEL_OTHER *regexp.Regexp
	LABEL_ALL   []*regexp.Regexp
)

func init() {
	LABEL_OWNER = regexp.MustCompile(`(owner\.)([_a-zA-Z0-9])+`)
	LABEL_OTHER = regexp.MustCompile(`(other\.)([_a-zA-Z0-9])+`)
	LABEL_ALL = []*regexp.Regexp{LABEL_OWNER, LABEL_OTHER}
}

func delRLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	loginName := r.Header.Get("User")
	if loginName == "" {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}
	r.ParseForm()

	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	var rep repository
	var err error
	Q := bson.M{COL_REPNAME: repname}
	if rep, err = db.getRepository(Q); err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}

	code, e := db.optLabels(r, LABEL_ALL, rep, Q, LABEL_OPT_SUB)
	return rsp.Json(code, e)
}

func upsertRLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	loginName := r.Header.Get("User")
	if loginName == "" {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}
	r.ParseForm()

	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	var rep repository
	var err error
	Q := bson.M{COL_REPNAME: repname}
	if rep, err = db.getRepository(Q); err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrRepositoryNotFound(repname))
	}

	code, e := db.optLabels(r, LABEL_ALL, rep, Q, LABEL_OPT_ADD)
	return rsp.Json(code, e)
}

func upsertDLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	loginName := r.Header.Get("User")
	if loginName == "" {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}
	r.ParseForm()

	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	var item dataItem
	var err error
	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	if item, err = db.getDataitem(Q); err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	code, e := db.optLabels(r, LABEL_ALL, item, Q, LABEL_OPT_ADD)

	return rsp.Json(code, e)
}

func delDLabelHandler(r *http.Request, rsp *Rsp, param martini.Params, db *DB) (int, string) {
	loginName := r.Header.Get("User")
	if loginName == "" {
		return rsp.Json(400, E(ErrorCodePermissionDenied))
	}
	r.ParseForm()

	repname := param["repname"]
	if repname == "" {
		return rsp.Json(400, ErrNoParameter("repname"))
	}

	itemname := param["itemname"]
	if itemname == "" {
		return rsp.Json(400, ErrNoParameter("itemname"))
	}

	var item dataItem
	var err error
	Q := bson.M{COL_REPNAME: repname, COL_ITEM_NAME: itemname}
	if item, err = db.getDataitem(Q); err != nil && err == mgo.ErrNotFound {
		return rsp.Json(400, ErrDataitemNotFound(itemname))
	}

	code, e := db.optLabels(r, LABEL_ALL, item, Q, LABEL_OPT_SUB)

	return rsp.Json(code, e)
}

func parseLabelParam(r *http.Request, reg *regexp.Regexp) []string {
	l := []string{}
	for k := range r.Form {
		str := reg.FindAllString(k, -1)
		if len(str) > 0 {
			l = append(l, str[0])
		}
	}
	return l
}

func (db *DB) optLabels(r *http.Request, reg []*regexp.Regexp, s interface{}, Q bson.M, optTp string) (int, *Error) {
	loginName := r.Header.Get("User")
	v := reflect.ValueOf(s)
	userName := v.FieldByName("Create_user").String()

	collectionName := C_DATAITEM
	if len(Q) == 1 {
		collectionName = C_REPOSITORY
	}

	for _, v := range reg {
		switch v {
		case LABEL_OWNER:
			if owner := parseLabelParam(r, LABEL_OWNER); len(owner) > 0 {
				if userName != loginName {
					return 400, E(ErrorCodePermissionDenied)
				}

				for _, v := range owner {
					colName := fmt.Sprintf("label.%s", v)
					values := r.Form[v]
					colValue := "1"
					switch optTp {
					case LABEL_OPT_ADD:
						if len(values) == 0 {
							return 400, ErrInvalidParameter(colName)
						}
						colValue = values[0]
					}

					exec := Execute{
						Collection: collectionName,
						Selector:   Q,
						Update:     bson.M{optTp: bson.M{colName: colValue}},
						Type:       Exec_Type_Update,
					}

					go asynExec(exec)
				}
			}
		case LABEL_OTHER:
			if other := parseLabelParam(r, LABEL_OTHER); len(other) > 0 {
				if getUserType(r, db) != USER_TP_ADMIN {
					return 400, E(ErrorCodePermissionDenied)
				}

				for _, v := range other {
					colName := fmt.Sprintf("label.%s", v)
					values := r.Form[v]
					colValue := "1"
					switch optTp {
					case LABEL_OPT_ADD:
						if len(values) == 0 {
							return 400, ErrInvalidParameter(colName)
						}
						colValue = values[0]
					}

					exec := Execute{
						Collection: collectionName,
						Selector:   Q,
						Update:     bson.M{optTp: bson.M{colName: colValue}},
						Type:       Exec_Type_Update,
					}

					go asynExec(exec)
				}
			}
		}
	}

	return 200, E(OK)
}
