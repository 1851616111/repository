package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	C_REPOSITORY            = "repository"
	C_REPOSITORY_DEL        = "repository_del"
	C_DATAITEM              = "dataitem"
	C_DATAITEM_DEL          = "dataitem_del"
	C_REPOSITORY_PERMISSION = "permission_rep"
	C_DATAITEM_PERMISSION   = "permission_item"
	C_SELECT                = "select"
	C_TAG                   = "tag"
	C_STATIS_DAY            = "statis_day"
)

type Label struct {
	Sys   map[string]interface{} `json:"sys"`
	Opt   map[string]interface{} `json:"opt"`
	Owner map[string]interface{} `json:"owner"`
	Other map[string]interface{} `json:"other"`
}

type repository struct {
	Repository_name string      `json:"-"`
	Create_user     string      `json:"create_user,omitempty"`
	Repaccesstype   string      `json:"repaccesstype,omitempty"`
	Comment         string      `json:"comment"`
	Optime          string      `json:"optime,omitempty"`
	Items           int         `json:"items"`
	CooperateItems  int         `json:"cooperateitems"`
	Label           interface{} `json:"label"`
	Ct              time.Time   `json:"-"`
	St              time.Time   `bson:"st,omitempty", json:"-"`
	Cooperate       interface{} `bson:"cooperators", json:"-"`
	Rank            float32     `json:"-"`
}

type Namelist []names
type names struct {
	Repository_name  string `json:"repname"`
	Cooperate_status string `json:"cooperatestate,omitempty"`
	Dataitem_name    string `json:"itemname,omitempty"`
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
	Cooperate       bool    `json:"-"`
	Rank            float32 `json:"-"`
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

func refreshDB() {

	for {
		select {
		case <-time.After(time.Second * 5):
			if err := db.Ping(); err != nil {
				Log.Infof("%s db connect err %s", time.Now().Format("2006-01-02 15:04:05"), err)
				db = DB{*connect()}
				db.Refresh()
			}
		}
	}
}

func connect() *mgo.Session {
	ip, port := getMgoAddr()
	if ip == "" {
		Log.Error("can not init mongo ip")
	}

	if port == "" {
		Log.Error("can not init mongo port")
	}

	url := fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, ip, port)
	Log.Infof("[Mongo Addr] %s", url)

	var session *mgo.Session
	var err error
	try := 0
	for {
		ip, port = getMgoAddr()
		url = fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, ip, port)
		session, err = mgo.Dial(url)
		if err != nil {
			try++
			Log.Errorf("dial mgo(%s) err %s, already try %d times", url, err.Error(), try)
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}

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
	TimeId          string      `json:"tid"`
	Type            string      `json:"type"`
	Repository_name interface{} `json:"repname"`
	Dataitem_name   interface{} `json:"itemname"`
	Tag             interface{} `json:"tag"`
	Time            string      `json:"time"`
}

type m_item struct {
	TimeId          string      `json:"tid"`
	Type            string      `json:"type"`
	Repository_name interface{} `json:"repname"`
	Dataitem_name   interface{} `json:"itemname"`
	Time            string      `json:"time"`
}

type m_rep struct {
	TimeId          string      `json:"tid"`
	Type            string      `json:"type"`
	Repository_name interface{} `json:"repname"`
	Time            string      `json:"time"`
}

func mqPermissionHandler(m map[string]interface{}) {

	if m[COL_REPNAME] != nil && m[COL_ITEM_NAME] != nil && m[COL_PERMIT_USER] != nil {
		exec := Execute{
			Collection: C_DATAITEM_PERMISSION,
			Update:     m,
			Type:       Exec_Type_Insert,
		}

		go asynExec(exec)
	}

}

func mqRankHandler(i interface{}) {
	execs := []Execute{}

	if list, ok := i.([]statisRepRank); ok {
		for _, statis := range list {
			exec := Execute{
				Collection: C_REPOSITORY,
				Selector:   bson.M{COL_REPNAME: statis.Repository_name},
				Update:     bson.M{CMD_SET: bson.M{COL_RANK: statis.Rank}},
				Type:       Exec_Type_Update,
			}
			execs = append(execs, exec)
		}
	}

	if list, ok := i.([]statisItemRank); ok {
		for _, statis := range list {
			exec := Execute{
				Collection: C_DATAITEM,
				Selector:   bson.M{COL_REPNAME: statis.Repository_name, COL_ITEM_NAME: statis.Dataitem_name},
				Update:     bson.M{CMD_SET: bson.M{COL_RANK: statis.Rank}},
				Type:       Exec_Type_Update,
			}
			execs = append(execs, exec)
		}
	}

	go asynExec(execs...)
}
