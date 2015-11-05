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

var (
	REG_BASIC_AUTH = regexp.MustCompile(`^Basic (.+)$`)
)

func auth(w http.ResponseWriter, r *http.Request, c martini.Context, db *DB) {
	email, passwd, err := parseBasicAuth(r)
	if err != nil {
		log.Println("basic parse err: ", err)
		if len(r.Header["X-Requested-With"]) == 0 {
			//send res header for non-xhr
			w.Header().Set("WWW-Authenticate", `Basic realm=Authorization Required"`)
		}
		http.Error(w, "unauthorized", 401)
		return
	}

	if exist, user_id := db.userExist(email, passwd); exist {
		c.Map(user_id)
		return
	}

	http.Error(w, "unauthorized", 401)
}

func parseBasicAuth(r *http.Request) (string, string, error) {
	s := r.Header.Get("Authorization")

	match := REG_BASIC_AUTH.FindAllStringSubmatch(s, -1)
	if match == nil {
		return "", "", fmt.Errorf("bad auth header %s\n", s)
	}

	s = match[0][1]

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", "", err
	}

	s = string(b)

	ary := strings.Split(s, ":")
	if len(ary) != 2 {
		return "", "", fmt.Errorf("bad auth string %s\n", s)
	}

	return ary[0], ary[1], nil
}
