package main

import (
	"fmt"
	"github.com/asiainfoLDP/datahub_repository/log"
	"github.com/asiainfoLDP/datahub_repository/mq"
	"github.com/go-martini/martini"
	"net/http"
)

var (
	DB_NAMESPACE_MONGO = "datahub"
	DB_NAME            = "datahub"
	SERVICE_PORT       = Env("goservice_port", false)

	Service_Name_Kafka = Env("kafka_service_name", true)
	Service_Name_Mongo = Env("mongo_service_name", true)

	DISCOVERY_CONSUL_SERVER_ADDR = Env("CONSUL_SERVER", true)
	DISCOVERY_CONSUL_SERVER_PORT = Env("CONSUL_DNS_PORT", true)

	db  DB    = DB{*connect()}
	q_c Queue = Queue{queue}
	msg Msg
	Log = log.NewLogger("http handler")
)

func main() {

	correctQuota(&db)
	initMq(getKFKAddr)

	go refreshDB()

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

		r.Delete("/repository/:repname/whitelist/:username", chkRepPermission, delRepPmsHandler)
		r.Delete("/repository/:repname/cooperator/:username", chkRepPermission, delRepCoptPmsHandler)
		r.Delete("/repository/:repname/dataitem/:itemname", chkItemPermission, delItemPmsHandler)
	})

	http.Handle("/", m)

	Log.Infof("service listen on %s", SERVICE_PORT)
	Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", SERVICE_PORT), nil))

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

	_, _, err = msg.SendSyncMessage(MQ_TOPIC_TO_REP, []byte(MQ_KEY_ADD_STATIS_RANK_REP), []byte("")) // force create the topic
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
	//return "10.1.235.98", "27017"
}

func getKFKAddr() (string, string) {
	entryList := dnsExchange(Service_Name_Kafka, DISCOVERY_CONSUL_SERVER_ADDR, DISCOVERY_CONSUL_SERVER_PORT)

	for _, v := range entryList {
		if v.port != "2181" {
			return v.ip, v.port
		}
	}
	return "", ""
	//return "10.1.235.98", "9092"
}
