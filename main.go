package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/quexer/utee"
	"net/http"
)

var (
	DB_NAMESPACE_MONGO = "datahub"
	DB_NAME            = "datahub"
	SERVICE_PORT       = utee.Env("goservice_port", false)

	DB_MONGO_USER   = utee.Env("DB_MONGO_USER", false)
	DB_MONGO_PASSWD = utee.Env("DB_MONGO_PASSWD", false)
	DB_MONGO_URL    = utee.Env("DB_MONGO_URL", false)
	DB_MONGO_PORT   = utee.Env("DB_MONGO_PORT", false)

	DB_URL_MONGO = fmt.Sprintf(`%s:%s/datahub?maxPoolSize=50`, DB_MONGO_URL, DB_MONGO_PORT)
	db           DB
	q_c          Queue
)

func init() {
	se := connect(DB_URL_MONGO)
	db = DB{*se}
	q_c = Queue{queueChannel}
}

func main() {

	go q_c.serve(&db)
	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		rsp := &Rsp{w: w}
		c.Map(rsp)
		c.Map(&db)
	})

	m.Post("/search", searchHandler)

	m.Group("/repositories", func(r martini.Router) {
		r.Get("", auth, getRsHandler)

		r.Get("/:repname", getRHandler)
		r.Get("/:repname/:itemname", getDHandler)
		r.Get("/:repname/:itemname/:tag", getTagHandler)

		r.Post("/:repname", auth, createRHandler)
		r.Post("/:repname/:itemname", auth, createDHandler)
		r.Post("/:repname/:itemname/:tag", auth, setTagHandler)

		r.Put("/:repname/:itemname", auth, updateDHandler)

		r.Put("/:repname", auth, getRsHandler)
		r.Delete("/:repname", auth, delRHandler)
		r.Delete("/:repname/:itemname", auth, delDHandler)
		r.Delete("/:repname/:itemname/:tag", auth, delTagHandler)
	})

	m.Group("/selects", func(r martini.Router) {
		r.Get("", getSelectsHandler)
		r.Post("", authAdmin, updateLabelHandler)
	})

	m.Group("/select_labels", func(r martini.Router) {
		r.Get("", getSelectLabelsHandler)
		r.Post("/:labelname", authAdmin, setSelectLabelHandler)
	})

	m.Group("/permit", func(r martini.Router) {
		r.Get("/:user_name", getUsrPmtRepsHandler)
		r.Post("/:user_name", setUsrPmtRepsHandler)
	}, auth)

	http.Handle("/", m)

	http.ListenAndServe(fmt.Sprintf(":%s", SERVICE_PORT), nil)

}
