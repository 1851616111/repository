package main

import (
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/quexer/utee"
	"net/http"
)

var (
	SERVICE_PORT = utee.Env("goservice_port", false)
	DB_ADDR      = utee.Env("MYSQL_PORT_3306_TCP_ADDR", false)
	DB_PORT      = utee.Env("MYSQL_PORT_3306_TCP_PORT", false)
	DB_DATABASE  = utee.Env("MYSQL_ENV_MYSQL_DATABASE", false)
	DB_USER      = utee.Env("MYSQL_ENV_MYSQL_USER", false)
	DB_PASSWORD  = utee.Env("MYSQL_ENV_MYSQL_PASSWORD", false)
	DB_URL       = fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8`, DB_USER, DB_PASSWORD, DB_ADDR, DB_PORT, DB_DATABASE)
	db           DB
	dim          Dim
)

func init() {

	engine, err := xorm.NewEngine("mysql", DB_URL)
	utee.Chk(err)
	db = DB{*engine}
	dim = Dim{mm: make(MM)}
}

func main() {

	go DimLoop(&db)

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		rsp := &Rsp{w: w}
		c.Map(rsp)
		c.Map(&db)
	})

	m.Group("/subscriptions", func(r martini.Router) {
		r.Get("", getSHandler)
		r.Get("/login", login)

	}, auth)
	m.Group("/repositories", func(r martini.Router) {
		r.Get("", getRHandler)
		r.Post("/:repname/:itemname", setDHandler)
	}, auth)

	m.Group("/repository", func(r martini.Router) {
		r.Get("/:repname/items", getRepoByNameHandler)
	})

	m.Group("/portal/dataitem", func(r martini.Router) {
		r.Get("", getItemsHandler)
		r.Post("/chosen", setItemChoseHandler)
		r.Get("/chosen", getItemChoseHandler)
		r.Get("/chosen/names", getChosenNamesHandler)
		//		r.Get("/test", test)
	})

	http.Handle("/", m)

	http.ListenAndServe(fmt.Sprintf(":%s", SERVICE_PORT), nil)

}
