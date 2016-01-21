package main

import (
	"fmt"
	"github.com/asiainfoLDP/datahub_repository/log"
	"github.com/asiainfoLDP/datahub_repository/mq"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

var (
	DB_NAMESPACE_MONGO = "datahub"
	DB_NAME            = "datahub"
	SERVICE_PORT       = Env("goservice_port", false)

	Service_Name_Kafka = Env("kafka_service_name", false)
	Service_Name_Mongo = "datahub_repository_mongo"

	DISCOVERY_CONSUL_SERVER_ADDR = Env("CONSUL_SERVER", false)
	DISCOVERY_CONSUL_SERVER_PORT = Env("CONSUL_DNS_PORT", false)

	MQ_KAFKA_ADDR = Env("MQ_KAFKA_ADDR", false)
	MQ_KAFKA_PORT = Env("MQ_KAFKA_PORT", false)

	db  DB = initDB()
	q_c Queue
	msg Msg
	Log = log.NewLogger("http handler")
)

func init() {
	if DISCOVERY_CONSUL_SERVER_ADDR == "" || DISCOVERY_CONSUL_SERVER_PORT == "" {
		Log.Fatal("can not get env CONSUL_SERVER CONSUL_DNS_PORT")
		os.Exit(0)
	}
	q_c = Queue{queue}
}

func main() {
	correctQuota(&db)
	initMq()

	//	go refreshDB(&db, func(db *DB) {
	//		ip, port := dnsExchange(Service_Name_Mongo, DISCOVERY_CONSUL_SERVER_ADDR, DISCOVERY_CONSUL_SERVER_PORT)
	//		db.Session = *getMgoSession(ip, port)
	//		db.Refresh()
	//	})

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

func initDB() DB {
	fmt.Println("func initDB() ", Service_Name_Mongo, DISCOVERY_CONSUL_SERVER_ADDR, DISCOVERY_CONSUL_SERVER_PORT)
	ip, port := dnsExchange(Service_Name_Mongo, DISCOVERY_CONSUL_SERVER_ADDR, DISCOVERY_CONSUL_SERVER_PORT)
	if ip == "" {
		Log.Error("------> mongo ip", ip)
		os.Exit(0)
	}

	if port == "" {
		Log.Error("------> mongo port", port)
		os.Exit(0)
	}
	return DB{*getMgoSession(ip, port)}
}

func initMq() {

	if MQ_KAFKA_ADDR == "" || MQ_KAFKA_PORT == "" {
		MQ_KAFKA_ADDR = "10.1.235.98"
		MQ_KAFKA_PORT = "9092"
	}

	MQ := fmt.Sprintf("%s:%s", MQ_KAFKA_ADDR, MQ_KAFKA_PORT)
	Log.Info(MQ)
	m_q, err := mq.NewMQ([]string{MQ})
	if err != nil {
		Log.Errorf("initMQ error: %s", err.Error())
		return
	}

	msg = Msg{m_q}

	myListener := newMyMesssageListener(MQ_HANDLER_PERMISSION)

	_, _, err = msg.SendSyncMessage(MQ_TOPIC_TO_REP, []byte(MQ_KEY_ADD_PERMISSION), []byte("")) // force create the topic
	get(err)

	err = msg.SetMessageListener(MQ_TOPIC_TO_REP, 0, mq.Offset_Marked, myListener)
	if err != nil {
		Log.Info("SetMessageListener error: ", err)
	}

}

func getMgoSession(mgoAddr, mgoPort string) *mgo.Session {
	DB_URL := fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, mgoAddr, mgoPort)
	Log.Info(DB_URL)
	return connect(DB_URL)
}
