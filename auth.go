package main

import (
	"github.com/go-martini/martini"
	"net/http"
)

const (
	ADMIN = "admin"
	PANXY = "panxy3@asiainfo.com"
)

func auth(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	if login_Name == "" {
		http.Error(w, "unauthorized", 401)
	}
	c.Map(login_Name)
	return

}

func authAdmin(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("User")
	if login_Name != ADMIN && login_Name != PANXY {
		http.Error(w, "unauthorized", 401)
		return
	}
	c.Map(login_Name)
	return
}
