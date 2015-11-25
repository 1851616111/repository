package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
)

var (
	DB_NAMESPACE_MONGO = "datahub"
	DB_NAME            = "datahub"
	SERVICE_PORT       = Env("goservice_port", false)

	DB_MONGO_USER   = Env("DB_MONGO_USER", false)
	DB_MONGO_PASSWD = Env("DB_MONGO_PASSWD", false)
	DB_MONGO_URL    = Env("MONGO_PORT_27017_TCP_ADDR", false)
	DB_MONGO_PORT   = Env("MONGO_PORT_27017_TCP_PORT", false)

	db  DB
	q_c Queue
	Log = NewLogger("http handler")
)

func init() {
	if DB_MONGO_URL == "" || DB_MONGO_PORT == "" {
		DB_MONGO_URL = "10.1.235.98"
		DB_MONGO_PORT = "27017"
	}

	DB_URL := fmt.Sprintf(`%s:%s/datahub?maxPoolSize=500`, DB_MONGO_URL, DB_MONGO_PORT)

	se := connect(DB_URL)
	db = DB{*se}
	q_c = Queue{queueChannel}
}

func main() {

	go q_c.serve(&db)
	go staticLoop(&db)
	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		rsp := &Rsp{w: w}
		c.Map(rsp)
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

		r.Post("/:repname", auth, createRHandler)
		r.Post("/:repname/:itemname", auth, createDHandler)
		r.Post("/:repname/:itemname/:tag", auth, setTagHandler)

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
		r.Delete("/:repname/:itemname", authAdmin, deleteSelectLabelHandler)
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

		r.Put("/:repname", chkRepPermission, upsertRepPmsHandler)
		r.Put("/:repname/:itemname", chkItemPermission, setItemPmsHandler)

		r.Delete("/:repname", chkRepPermission, delRepPmsHandler)
		r.Delete("/:repname/:itemname", chkItemPermission, delItemPmsHandler)
	})

	http.Handle("/", m)

	Log.Infof("service listen on %s", SERVICE_PORT)
	Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", SERVICE_PORT), nil))

}
