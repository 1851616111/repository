package main

import (
	"github.com/quexer/utee"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	C_REPOSITORY        = "repository"
	C_DATAITEM          = "dataitem"
	C_REPOSITORY_PERMIT = "permituser1"
	C_SELECT            = "select"
	C_TAG               = "tag"
)

type label struct {
	Sys   struct{} `json:"sys"`
	Opt   struct{} `json:"opt"`
	Owner struct{} `json:"owner"`
	Other struct{} `json:"other"`
}

type repository struct {
	Repository_name string    `json:"-"`
	Create_user     string    `json:"create_user,omitempty"`
	Repaccesstype   string    `json:"repaccesstype,omitempty"`
	Deposit         bool      `json:"deposit"`
	Comment         string    `json:"comment"`
	Optime          time.Time `json:"optime,omitempty"`
	Stars           int       `json:"stars"`
	Views           int       `json:"views"`
	Items           int       `json:"items"`
	Label           *label    `json:"label,omitempty"`
	Ct              time.Time `json:"-"`
	St              time.Time `json:"-"`
}
type names struct {
	Repository_name string `json:"repname"`
	Dataitem_name   string `json:"itemname"`
}
type dataItem struct {
	Repository_name string    `json:"-"`
	Create_name     string    `json:"create_user,omitempty"`
	Dataitem_name   string    `json:"-"`
	Itemaccesstype  string    `json:"itemaccesstype,omitempty"`
	Price           string    `json:"price,omitempty"`
	Optime          time.Time `json:"optime,omitempty"`
	Meta            string    `json:"meta"`
	Sample          string    `json:"sample"`
	Comment         string    `json:"comment"`
	Label           label     `json:"label"`
	Ct              time.Time `json:"-"`
	St              time.Time `json:"-"`
}

type Tag struct {
	Dataitem_id int64  `json:"dataitem_id,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Optime      string `json:"optime,omitempty"`
}

type Repository_Permit struct {
	User_name       string `json:"-"`
	Repository_name string `json:"repository_name"`
}

type Dim_Table struct {
	Field_name string
	Id         int
	Name       string
}

type Select struct {
	LabelName string `json:"labelname,omitempty" `
	Order     int    `json:"-"`
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
	err := db.C(C_REPOSITORY).EnsureIndex(mgo.Index{Key: []string{COL_REP_NAME}, Unique: true})
	utee.Chk(err)
	err = db.C(C_DATAITEM).EnsureIndex(mgo.Index{Key: []string{COL_REP_NAME, COL_ITEM_NAME}, Unique: true})
	utee.Chk(err)
	err = db.C(C_SELECT).EnsureIndex(mgo.Index{Key: []string{COL_SELECT_LABEL}, Unique: true})
	utee.Chk(err)
	err = db.C(C_REPOSITORY_PERMIT).EnsureIndex(mgo.Index{Key: []string{COL_REP_NAME, COL_PERMIT_USER}, Unique: true})
	utee.Chk(err)

}
