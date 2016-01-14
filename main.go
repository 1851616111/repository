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

	DB_MONGO_ADDR = Env("MONGO_PORT_27017_TCP_ADDR", false)
	DB_MONGO_PORT = Env("MONGO_PORT_27017_TCP_PORT", false)
	MQ_KAFKA_ADDR = Env("MQ_KAFKA_ADDR", false)
	MQ_KAFKA_PORT = Env("MQ_KAFKA_PORT", false)

	db  DB = initDB()
	q_c Queue
	msg Msg
	Log = log.NewLogger("http handler")
)

func init() {
	q_c = Queue{queue}
}

func main() {

	initMq()

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

	if DB_MONGO_ADDR == "" || DB_MONGO_PORT == "" {
		DB_MONGO_ADDR = "10.1.235.98"
		DB_MONGO_PORT = "27017"
	}

	DB_URL := fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, DB_MONGO_ADDR, DB_MONGO_PORT)
	Log.Info(DB_URL)
	return DB{*connect(DB_URL)}

}

func initMq() {

	if MQ_KAFKA_ADDR == "" || MQ_KAFKA_PORT == "" {
		MQ_KAFKA_ADDR = DB_MONGO_ADDR
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
