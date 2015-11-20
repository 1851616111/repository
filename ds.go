package main

import (
	"github.com/quexer/utee"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	C_REPOSITORY            = "repository"
	C_DATAITEM              = "dataitem"
	C_REPOSITORY_PERMISSION = "permission_rep"
	C_DATAITEM_PERMISSION   = "permission_item"
	C_SELECT                = "select"
	C_TAG                   = "tag"
)

//type label struct {
//	Sys   interface{} `json:"sys"`
//	Opt   interface{} `json:"opt"`
//	Owner interface{} `json:"owner"`
//	Other interface{} `json:"other"`
//}

type repository struct {
	Repository_name string `json:"-"`
	Create_user     string `json:"create_user,omitempty"`
	Repaccesstype   string `json:"repaccesstype,omitempty"`
	//	Deposit         bool      `json:"deposit"`
	Comment string      `json:"comment"`
	Optime  string      `json:"optime,omitempty"`
	Items   int         `json:"items"`
	Label   interface{} `json:"label"`
	Ct      time.Time   `json:"-"`
	st      time.Time
}
type names struct {
	Repository_name string `json:"repname"`
	Dataitem_name   string `json:"itemname,omitempty"`
}
type search struct {
	Repository_name string
	Dataitem_name   string
	Ct              time.Time
}
type dataItem struct {
	Repository_name string      `json:"-"`
	Dataitem_name   string      `json:"-"`
	Create_user     string      `json:"create_user,omitempty"`
	Itemaccesstype  string      `json:"itemaccesstype,omitempty"`
	Price           interface{} `bson:"-", json:"price,omitempty"`
	Optime          string      `json:"optime,omitempty"`
	Meta            string      `bson:"-", json:"meta"`
	Sample          string      `bson:"-", json:"sample"`
	Comment         string      `json:"comment"`
	Tags            int         `json:"tags"`
	Label           interface{} `json:"label"`
	Ct              time.Time   `json:"-"`
	st              time.Time
}

type tag struct {
	Repository_name string `json:"-"`
	Dataitem_name   string `json:"-"`
	Tag             string `json:"tag,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Optime          string `json:"optime,omitempty"`
}

type Rep_Permission struct {
	User_name       string `json:"username"`
	Repository_name string `json:"repname"`
	Write           int    `bson:",omitempty", json:"write"`
}

type Item_Permission struct {
	User_name     string `json:"username"`
	Dataitem_name string `json:"itemname"`
}

type Dim_Table struct {
	Field_name string
	Id         int
	Name       string
}

type Select struct {
	LabelName    string `json:"labelname,omitempty"`
	NewLabelName string `json:"newLabelName,omitempty" bson:"-"`
	Order        int    `json:"order,omitempty"`
	Icon         string `json:"icon,omitempty"`
}

type DB struct {
	mgo.Session
}

func connect(db_connection string) *mgo.Session {
	session, err := mgo.Dial(db_connection)
	utee.Chk(err)
	initDb(session)
	return session
}

//初始化索引
func initDb(session *mgo.Session) {
	db := session.DB(DB_NAMESPACE_MONGO)
	err := db.C(C_REPOSITORY).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME}, Unique: true})
	utee.Chk(err)
	err = db.C(C_DATAITEM).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_ITEM_NAME}, Unique: true})
	utee.Chk(err)
	err = db.C(C_SELECT).EnsureIndex(mgo.Index{Key: []string{COL_SELECT_LABEL}, Unique: true})
	utee.Chk(err)
	err = db.C(C_REPOSITORY_PERMISSION).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_PERMIT_USER}, Unique: true})
	utee.Chk(err)
	err = db.C(C_DATAITEM_PERMISSION).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_PERMIT_USER}, Unique: true})
	utee.Chk(err)
	err = db.C(C_TAG).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_ITEM_NAME, COL_TAG_NAME}, Unique: true})
	utee.Chk(err)
	err = db.C(C_DATAITEM).EnsureIndexKey(COL_REPNAME)
	utee.Chk(err)
	err = db.C(C_DATAITEM).EnsureIndexKey(COL_ITEM_NAME)
	utee.Chk(err)
	err = db.C(C_DATAITEM).EnsureIndexKey(COL_COMMENT)
	utee.Chk(err)
}
