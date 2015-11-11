package main

import (
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/quexer/utee"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	INNER = "/inner"
)

var (
	NAMESPACE    = "datahub"
	DB_NAME      = "datahub"
	SERVICE_PORT = utee.Env("goservice_port", false)

	MONGO_ADDR = utee.Env("DB_MONGO_URL", false)
	MONGO_PORT = utee.Env("DB_MONGO_PORT", false)

	MONGO_URL = fmt.Sprintf(`%s:%s/datahub?maxPoolSize=50`, MONGO_ADDR, MONGO_PORT)
	db        DB
	m_db      *mgo.Session
)

func init() {

	m_db = connect(MONGO_URL)
	db = DB{*m_db}
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
