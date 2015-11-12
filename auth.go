package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	ADMIN = "admin"
)

var (
	REG_BASIC_AUTH = regexp.MustCompile(`^Basic (.+)$`)
)

func auth(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("user")
	c.Map(login_Name)
	return

}

func authAdmin(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	login_Name := r.Header.Get("user")
	if login_Name != ADMIN {
		http.Error(w, "unauthorized", 401)
		return
	}
	c.Map(login_Name)
	return
}
