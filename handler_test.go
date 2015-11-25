package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type param struct {
	requestBody string
	rsp         *Rsp
	param       martini.Params
	db          *DB
	login_name  string
}

type Body struct {
	Result
}

type expect struct {
	code int
	body Body
}
type Context struct {
	description string
	param       param
	expect      expect
}

var (
	ramdom int
)

func init() {
	rd := rand.New(rand.NewSource(int64(time.Now().Second())))
	ramdom = rd.Int()
}

func Test_createRHandler(t *testing.T) {

	contexts := []Context{
		Context{
			description: "----------> create repository",
			param: param{
				requestBody: `{
									"repaccesstype": "public",
									"comment": "中国移动北京终端详情",
									"label": {
										"sys": {
											"loc": "北京"
											},
										"opt": {
											"age": 22
											},
										"owner": {
											"name": "michael"
											},
										"other": {
											"friend": 22
										}
									}
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom)},
				db:         &db,
				login_name: "panxy3@asiainfo.com",
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
	}

	for _, v := range contexts {
		t.Log(v.description)
		p := v.param
		r, err := http.NewRequest("POST", "/repositories/rep0001", strings.NewReader(p.requestBody))
		get(err)

		code, msg := createRHandler(r, p.rsp, p.param, p.db, p.login_name)
		//		res := Result{Code: code, Msg: msg}
		//		v.expected.expect(t, res)
		log.Println(code)
		log.Println(msg)

	}
}

func (expect *expect) expect(t testing.T, resutlCode int, resutlData string) bool {
	if expect.code != resutlCode {
		t.Errorf("expected http.code:%d != return http.code:%d", expect.code, expect.code)
		return false
	}

	res := new(Result)
	err := json.Unmarshal([]byte(resutlData), res)
	t.Log(err)

	if expect.body.Code != res.Code {
		t.Errorf("expected http.Code(%d) != return http.Code(%d)", expect.code, expect.code)
		return false
	}

	if expect.body.Msg != res.Msg {
		t.Errorf("expected http.Msg(%d) != return http.Msg(%d)", expect.body.Msg, res.Msg)
		return false
	}

	if expect.body.Data != res.Data {
		t.Errorf("expected http.Data(%d) != return http.Data(%d)", expect.body.Data, res.Data)
		return false
	}

	return true
}
