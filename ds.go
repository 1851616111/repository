package main

import "github.com/go-xorm/xorm"

const (
	REPOSITORY        = "DH_REPOSITORY"
	DATAITEM          = "DH_DATAITEM"
	REPOSITORY_PERMIT = "DH_PERMITUSER1"
	DIM               = "DH_DIMTABLE"
	USER              = "DH_USER"
	DATAITEM_CHOSEN   = "DH_DATAITEM_CHOSEN"
	DATAITEMUSAGE     = "DH_DATAITEMUSAGE"
	TAG               = "DH_TAG"
)

type Repository struct {
	Repository_name string `json:"repository_name,omitempty"`
	User_id         int    `json:"user_id,omitempty"`
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

type DataItem struct {
	Repository_name string  `json:"repname,omitempty"`
	User_id         int     `json:"user_id,omitempty"`
	Dataitem_id     int64   `json:"dataitem_id,omitempty"  xorm:"dataitem_id pk autoincr"`
	Dataitem_name   string  `json:"dataitem_name,omitempty"`
	Ico_name        string  `json:"ico_name,omitempty"`
	Permit_type     int     `json:"permit_type,omitempty"`
	Key_words       string  `json:"key_words,omitempty"`
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

type Dataitem_Chosen struct {
	Chosen_name string `json:"chosen_name,omitempty"`
	Dataitem_id int    `json:"dataitem_id,omitempty"`
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
func (p *Dataitem_Chosen) TableName() string {
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
	xorm.Engine
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
