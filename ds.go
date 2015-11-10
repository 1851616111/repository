package main

import (
	"github.com/quexer/utee"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	REPOSITORY        = "DH_REPOSITORY"
	DATAITEM          = "DH_DATAITEM"
	REPOSITORY_PERMIT = "DH_PERMITUSER1"
	DIM               = "DH_DIMTABLE"
	USER              = "DH_USER"
	DATAITEM_CHOSEN   = "DH_DATAITEM_CHOSEN"
	DATAITEMUSAGE     = "DH_DATAITEMUSAGE"
	TAG               = "DH_TAG"

	M_REPOSITORY        = "repository"
	M_DATAITEM          = "dataitem"
	M_REPOSITORY_PERMIT = "PERMITUSER1"
	M_DIM               = "DIMTABLE"
	M_USER              = "USER"
	M_SELECT            = "SELECT"
	M_TAG               = "TAG"
)

type label struct {
	Sys   struct{} `json:"sys"`
	Opt   struct{} `json:"opt"`
	Owner struct{} `json:"owner"`
	Other struct{} `json:"other"`
}

type Repository struct {
	Repository_name string `json:"repository_name,omitempty"`
	Login_name      int    `json:"login_name,omitempty"`
	Permit_type     int    `json:"permit_type,omitempty"`
	Arrang_type     int    `json:"arrang_type,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Rank            int    `json:"rank,omitempty"`
	Status          int    `json:"status,omitempty"`
	Dataitems       int    `json:"dataitems,omitempty"`
	Tags            int    `json:"tags,omitempty"`
	Stars           int    `json:"stars,omitempty"`
	Optime          string `json:"optime,omitempty"`
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
	Label           *label    `json:"label,omitempty"`
}
type DataItem struct {
	Repository_name string  `json:"repname,omitempty"`
	Login_name      int     `json:"login_name,omitempty"`
	Dataitem_id     int64   `json:"dataitem_id,omitempty"  xorm:"dataitem_id pk autoincr"`
	Dataitem_name   string  `json:"dataitem_name,omitempty"`
	Ico_name        string  `json:"ico_name,omitempty"`
	Permit_type     int     `json:"permit_type,omitempty"`
	Label           string  `json:"label,omitempty"`
	Supply_style    int     `json:"supply_style,omitempty"`
	Priceunit_type  int     `json:"priceunit_type,omitempty"`
	Price           float32 `json:"price,omitempty"`
	Optime          string  `json:"optime,omitempty"`
	Data_format     int     `json:"data_format,omitempty"`
	Refresh_type    string  `json:"refresh_type,omitempty"`
	Refresh_num     int     `json:"refresh_num,omitempty"`
	Meta_filename   string  `json:"meta_filename,omitempty"`
	Sample_filename string  `json:"sample_filename,omitempty"`
	Comment         string  `json:"comment,omitempty"`
}

type DataItemUsage struct {
	Dataitem_id   int64  `json:"dataitem_id,omitempty"  xorm:"dataitem_id pk autoincr"`
	Dataitem_name string `json:"dataitem_name,omitempty"`
	Views         int    `json:"views"`
	Follows       int    `json:"follows"`
	//	Downloads     int    `json:"downloads"`
	Stars        int    `json:"stars"`
	Refresh_date string `json:"refresh_date,omitempty"`
	Usability    int    `json:"usability,omitempty"`
}

type Tag struct {
	Dataitem_id int64  `json:"dataitem_id,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Optime      string `json:"optime,omitempty"`
}

type User struct {
	User_id      int    `json:"user_id,omitempty"`
	User_status  int    `json:"user_status,omitempty"`
	User_type    int    `json:"user_type,omitempty"`
	Arrange_type int    `json:"arrange_type,omitempty"`
	Login_name   string `json:"login_name,omitempty"`
	Email        string `json:"email,omitempty"`
	Login_passwd string `json:"login_passwd,omitempty"`
	Sell_level   int    `json:"sell_level,omitempty"`
}

type Repository_Permit struct {
	Repository_name string
	User_id         int
}

type Dim_Table struct {
	Field_name string
	Id         int
	Name       string
}

type Select struct {
	LabelName       string `json:"labelname,omitempty" `
	Index           int    `json:"index,omitempty"`
	Dataitem_name   string `json:"dataitem_name,omitempty"`
	Repository_name string `json:"repository_name,omitempty"`
}

func (p *Repository) TableName() string {
	return REPOSITORY
}
func (p *DataItem) TableName() string {
	return DATAITEM
}
func (p *DataItemUsage) TableName() string {
	return DATAITEMUSAGE
}
func (p *Repository_Permit) TableName() string {
	return REPOSITORY_PERMIT
}
func (p *Dim_Table) TableName() string {
	return DIM
}
func (p *User) TableName() string {
	return USER
}
func (p *Select) TableName() string {
	return DATAITEM_CHOSEN
}
func (p *Tag) TableName() string {
	return TAG
}

type Data struct {
	Item  *DataItem      `json:"item,omitempty"`
	Usage *DataItemUsage `json:"statis,omitempty"`
	Tags  []Tag          `json:"tags,omitempty"`
}

type DB struct {
	mgo.Session
}

type Result struct {
	Repname       string     `json:"repname"`
	Repaccesstype string     `json:"repaccesstype"`
	Items         []DataItem `json:"items"`
}
type Subscription struct {
	Subscription_id int     `json:"subscription_id,omitempty"`
	User_id         int     `json:"user_id,omitempty"`
	Dataitem_id     int64   `json:"dataitem_id,omitempty"`
	Amount          int     `json:"amount,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Optime          string  `json:"optime,omitempty"`
}

func connect(db_connection string) *mgo.Session {
	session, err := mgo.Dial(db_connection)
	utee.Chk(err)
	initDb(session)
	return session
}

//初始化索引
func initDb(session *mgo.Session) {
	db := session.DB(NAMESPACE)
	err := db.C(M_REPOSITORY).EnsureIndex(mgo.Index{Key: []string{"repository_name"}, Unique: true})
	utee.Chk(err)
	err = db.C(M_DATAITEM).EnsureIndex(mgo.Index{Key: []string{"repository_name", "dataitem_name"}, Unique: true})
	utee.Chk(err)
	err = db.C(M_SELECT).EnsureIndex(mgo.Index{Key: []string{"repository_name", "dataitem_name"}, Unique: true})
	utee.Chk(err)

}
