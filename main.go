package main

import (
	"fmt"
	"github.com/asiainfoLDP/datahub_repository/log"
	"github.com/asiainfoLDP/datahub_repository/mq"
	"github.com/go-martini/martini"
	"net/http"
	"encoding/json"
)

var (
	DB_NAMESPACE_MONGO = "datahub"
	DB_NAME            = "datahub"
	SERVICE_PORT       = Env("goservice_port", false)

	Service_Name_Kafka = Env("kafka_service_name", false)
	Service_Name_Mongo = "datahub_repository_mongo"

	DISCOVERY_CONSUL_SERVER_ADDR = Env("CONSUL_SERVER", false)
	DISCOVERY_CONSUL_SERVER_PORT = Env("CONSUL_DNS_PORT", false)

	db  DB = initDB(getMgoAddr)
	q_c Queue
	msg Msg
	Log = log.NewLogger("http handler")
)

func init() {
	if DISCOVERY_CONSUL_SERVER_ADDR == "" || DISCOVERY_CONSUL_SERVER_PORT == "" {
		Log.Fatal("can not get env CONSUL_SERVER CONSUL_DNS_PORT")
	}

	if Service_Name_Kafka == "" {
		Log.Fatal("can not get env datahub_repository_mongo")
	}
	q_c = Queue{queue}

}

func main() {

	correctQuota(&db)
	initMq(getKFKAddr)

	go refreshDB(&db, func(db *DB) {
		ip, port := getMgoAddr()
		DB_URL := fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, ip, port)
		db.Session = *connect(DB_URL)
		db.Refresh()
	})

	go q_c.serve(&db)
	go staticLoop(&db)
	go pushMetaDataLoop(&db)

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		rsp := &Rsp{w: w}
		c.Map(rsp)
		c.Map(&msg)
		copy := db.Copy()
		c.Map(&DB{*copy})

	})

	m.Get("/search", searchHandler)

	m.Group("/repositories", func(r martini.Router) {
		r.Get("", getRsHandler)

		r.Get("/deleted", getDetetedHandler)
		r.Get("/:repname", getRHandler)
		r.Get("/:repname/:itemname", getDHandler)
		r.Get("/:repname/:itemname/subpermission", getDWithPermissionHandler)
		r.Get("/:repname/:itemname/:tag", getTagHandler)

		r.Post("/:repname", auth, getQuota, createRHandler)
		r.Post("/:repname/:itemname", auth, createDHandler)
		r.Post("/:repname/:itemname/:tag", auth, createTagHandler)

		r.Put("/:repname", auth, updateRHandler)
		r.Put("/:repname/label", upsertRLabelHandler)
		r.Put("/:repname/:itemname", auth, updateDHandler)
		r.Put("/:repname/:itemname/label", auth, upsertDLabelHandler)
		r.Put("/:repname/:itemname/:tag", auth, updateTagHandler)

		r.Delete("/:repname", auth, delRHandler)
		r.Delete("/:repname/label", delRLabelHandler)
		r.Delete("/:repname/:itemname", auth, delDHandler)
		r.Delete("/:repname/:itemname/label", delDLabelHandler)
		r.Delete("/:repname/:itemname/:tag", auth, delTagHandler)

	})

	m.Group("/selects", func(r martini.Router) {
		r.Get("", getSelectsHandler)
		r.Put("/:repname/:itemname", authAdmin, updateSelectHandler)
		r.Post("/:repname/:itemname", authAdmin, updateSelectHandler)
		r.Delete("/:repname/:itemname", authAdmin, delSelectHandler)
	})

	m.Group("/select_labels", func(r martini.Router) {
		r.Get("", getSelectLabelsHandler)
		r.Put("/:labelname", authAdmin, updateSelectLabelHandler)
		r.Post("/:labelname", authAdmin, setSelectLabelHandler)
		r.Delete("/:labelname", authAdmin, delSelectLabelHandler)
	})

	m.Group("/permit", func(r martini.Router) {
		r.Get("/:user_name", getUsrPmtRepsHandler)
		r.Post("/:user_name", setUsrPmtRepsHandler)
	}, auth)

	m.Group("/permission", func(r martini.Router) {
		r.Get("/:repname", chkRepPermission, getRepPmsHandler)
		r.Get("/:repname/:itemname", chkItemPermission, getItemPmsHandler)

		r.Put("/:repname", chkRepPermission, setRepPmsHandler)
		r.Put("/:repname/:itemname", chkItemPermission, setItemPmsHandler)

		r.Delete("/:repname", chkRepPermission, delRepPmsHandler)
		r.Delete("/:repname/:itemname", chkItemPermission, delItemPmsHandler)
	})

	http.Handle("/", m)

	Log.Infof("service listen on %s", SERVICE_PORT)
	Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", SERVICE_PORT), nil))

}

func initDB(f func() (string, string)) DB {
	ip, port := f()
	if ip == "" {
		Log.Error("can not init mongo ip")
	}

	if port == "" {
		Log.Error("can not init mongo port")
	}

	DB_URL := fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, ip, port)
	Log.Infof("[Mongo Addr] %s", DB_URL)

	return DB{*connect(DB_URL)}
}

func initMq(f func() (string, string)) {
	ip, port := f()
	if ip == "" {
		Log.Error("can not init kafka ip")
	}

	if port == "" {
		Log.Error("can not init kafka port")
	}

	MQ := fmt.Sprintf("%s:%s", ip, port)
	Log.Infof("[Kafka Addr] %s", MQ)
	m_q, err := mq.NewMQ([]string{MQ})
	if err != nil {
		Log.Errorf("initMQ error: %s", err.Error())
		return
	}

	msg = Msg{m_q}

	myListener := newMyMesssageListener(MQ_HANDLER_PERMISSION)

	r := new(repository)
	r.Create_user = "panxy3@asiainfo.com"
	r.Repository_name = "test_rank"
	b, _ := json.Marshal(r)
	_, _, err = msg.SendSyncMessage(MQ_TOPIC_TO_REP, []byte(MQ_KEY_ADD_STATIS_RANK), b) // force create the topic
	get(err)

	err = msg.SetMessageListener(MQ_TOPIC_TO_REP, 0, mq.Offset_Marked, myListener)
	if err != nil {
		Log.Info("SetMessageListener error: ", err)
	}

}

func getMgoAddr() (string, string) {
	entryList := dnsExchange(Service_Name_Mongo, DISCOVERY_CONSUL_SERVER_ADDR, DISCOVERY_CONSUL_SERVER_PORT)

	if len(entryList) > 0 {
		return entryList[0].ip, entryList[0].port
	}
	return "", ""
}

func getKFKAddr() (string, string) {
	entryList := dnsExchange(Service_Name_Kafka, DISCOVERY_CONSUL_SERVER_ADDR, DISCOVERY_CONSUL_SERVER_PORT)

	for _, v := range entryList {
		if v.port != "2181" {
			return v.ip, v.port
		}
	}
	return "", ""
}
