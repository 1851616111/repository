package main

import (
	"log"
	"net/http"
	"testing"
)

func Test_updateUser(t *testing.T) {
	l := Limit{100, 100}
	r := &http.Request{}

	r.Header.Set("User", "panxy3@asiainfo.com")
	r.Header.Set(AUTHORIZATION, "4d9e6e65942b31d861f5a5087c36b5b7")
	b, err := updateUser(r, l)
	t.Log((string(b)))
	t.Log(err)
	log.Println(string(b))
	log.Println(err)
}
