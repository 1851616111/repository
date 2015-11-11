package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/quexer/utee"
	"net/http"
)

var (
	DB_NAMESPACE_MONGO = "datahub"
	DB_NAME_MONGO      = "datahub"
	SERVICE_PORT       = utee.Env("goservice_port", false)

	DB_ADDR_MONGO = utee.Env("DB_MONGO_URL", false)
	DB_PORT_MONGO = utee.Env("DB_MONGO_PORT", false)

	DB_URL_MONGO = fmt.Sprintf(`%s:%s/datahub?maxPoolSize=50`, DB_ADDR_MONGO, DB_PORT_MONGO)
	db           DB
)

func init() {
	se := connect(DB_URL_MONGO)
	db = DB{*se}
}

func main() {

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		rsp := &Rsp{w: w}
		c.Map(rsp)
		c.Map(&db)
	})

	m.Group("/repositories", func(r martini.Router) {
		r.Get("", getRsHandler)

		r.Get("/:repname", getRHandler)
		r.Get("/:repname/:itemname", getDHandler)

		r.Post("/:repname", createRHandler)
		r.Post("/:repname/:itemname", createDHandler)
		r.Post("/:repname/:itemname/:tag", setTagHandler)

		r.Put("/:repname", getRsHandler)
		r.Delete("/:repname", getRsHandler)
	})

	m.Group("/selects", func(r martini.Router) {
		r.Get("", getSelectsHandler)
		r.Post("", updateLabelHandler)
	})

	m.Group("/select_labels", func(r martini.Router) {
		r.Get("", getSelectLabelsHandler)
		r.Post("/:labelname", setSelectLabelHandler)
	})

	http.Handle("/", m)

	http.ListenAndServe(fmt.Sprintf(":%s", SERVICE_PORT), nil)

}
