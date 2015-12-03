package main

import (
	"encoding/json"
	"github.com/asiainfoLDP/datahub_messages/mq"
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
	MQ_TOPIC                = "repositories_events_json"
)

var (
	COL_LABEL_CHILDREN = []string{"sys", "opt", "owner", "other"}
)

type Label struct {
	Sys   map[string]interface{} `json:"sys"`
	Opt   map[string]interface{} `json:"opt"`
	Owner map[string]interface{} `json:"owner"`
	Other map[string]interface{} `json:"other"`
}

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

type Namelist []names
type names struct {
	Repository_name string `json:"repname"`
	Dataitem_name   string `json:"itemname,omitempty"`
}
type search struct {
	Repository_name string
	Dataitem_name   string
	Optime          string
}
type dataItem struct {
	Repository_name string      `json:"-"`
	Dataitem_name   string      `json:"-"`
	Create_user     string      `json:"create_user,omitempty"`
	Itemaccesstype  string      `json:"itemaccesstype,omitempty"`
	Price           interface{} `json:"price,omitempty"`
	Optime          string      `json:"optime,omitempty"`
	Meta            string      `bson:"-", json:"meta"`
	Sample          string      `bson:"-", json:"sample"`
	Comment         string      `json:"comment"`
	Tags            int         `json:"tags"`
	Label           interface{} `json:"label"`
	Ct              time.Time   `json:"-"`
	st              time.Time
}

func (rep *repository) chkLabel() {
	if m, ok := rep.Label.(map[string]interface{}); ok {
		for _, v := range COL_LABEL_CHILDREN {
			if _, ok := m[v]; !ok {
				m[v] = make(map[string]interface{})
			}
		}
	} else {
		rep.Label = new(Label)
	}
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
	Repository_name string `json:"-"`
	Opt_permission  int    `json:"opt_permission"`
}

type Item_Permission struct {
	User_name       string `json:"username"`
	Repository_name string `json:"-"`
	Dataitem_name   string `json:"-"`
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

func (db *DB) copy() *DB {
	return &DB{*db.Copy()}
}
func connect(db_connection string) *mgo.Session {
	session, err := mgo.Dial(db_connection)
	get(err)
	initDb(session)
	return session
}

//初始化索引
func initDb(session *mgo.Session) {
	db := session.DB(DB_NAMESPACE_MONGO)
	err := db.C(C_REPOSITORY).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME}, Unique: true})
	get(err)
	err = db.C(C_DATAITEM).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_ITEM_NAME}, Unique: true})
	get(err)
	err = db.C(C_SELECT).EnsureIndex(mgo.Index{Key: []string{COL_SELECT_LABEL}, Unique: true})
	get(err)
	err = db.C(C_REPOSITORY_PERMISSION).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_PERMIT_USER}, Unique: true})
	get(err)
	err = db.C(C_DATAITEM_PERMISSION).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_ITEM_NAME, COL_PERMIT_USER}, Unique: true})
	get(err)
	err = db.C(C_TAG).EnsureIndex(mgo.Index{Key: []string{COL_REPNAME, COL_ITEM_NAME, COL_TAG_NAME}, Unique: true})
	get(err)
	err = db.C(C_DATAITEM).EnsureIndexKey(COL_REPNAME)
	get(err)
	err = db.C(C_DATAITEM).EnsureIndexKey(COL_ITEM_NAME)
	get(err)
	err = db.C(C_DATAITEM).EnsureIndexKey(COL_COMMENT)
	get(err)
}

type m_tag struct {
	Type            string      `json:"type"`
	Repository_name interface{} `json:"repname"`
	Dataitem_name   interface{} `json:"itemname"`
	Tag             interface{} `json:"tag"`
	Time            string      `json:"time"`
}

type m_item struct {
	Type            string      `json:"type"`
	Repository_name interface{} `json:"repname"`
	Dataitem_name   interface{} `json:"itemname"`
	Time            string      `json:"time"`
}

type m_rep struct {
	Type            string      `json:"type"`
	Repository_name interface{} `json:"repname"`
	Time            string      `json:"time"`
}

type Msg struct {
	mq.KafukaMQ
}

func (m *Msg) MqJson(content interface{}) {
	b, _ := json.Marshal(content)
	msg.SendAsyncMessage(MQ_TOPIC, []byte(""), b)
}

type MyMesssageListener struct {
	name string
}

func newMyMesssageListener(name string) *MyMesssageListener {
	return &MyMesssageListener{name: name}
}

func (listener *MyMesssageListener) OnMessage(key, value []byte, topic string, partition int32, offset int64) {
	Log.Errorf("%s received: (%d) message: %s", listener.name, offset, string(value))
}

func (listener *MyMesssageListener) OnError(err error) bool {
	Log.Errorf("api response listener error: %s", err.Error())
	return false
}
