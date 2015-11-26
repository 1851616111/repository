package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/asiainfoLDP/datahub_subscriptions/log"
	"github.com/go-martini/martini"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	USERNAME = "panxy3@asiainfo.com"
)

func errDatabaseOperate() string {
	return fmt.Sprintf("database operate : E11000 duplicate key error index: datahub.repository.$repository_name_1 dup key: { : \"app0001_%d\" }", ramdom)
}

type param struct {
	requestBody string
	rsp         *Rsp
	param       martini.Params
	db          *DB
	login_name  string
	repName     string
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
	rd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	ramdom = rd.Int()
}

func Test_createRHandler(t *testing.T) {

	contexts := []Context{
		Context{
			description: "Test_createRHandler ----------> create repository 1 ",
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
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_createRHandler ----------> create repository 2 ",
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
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1008,
					Msg:  errDatabaseOperate(),
				}},
			},
		},
		Context{
			description: "Test_createRHandler ----------> create repository 3 ",
			param: param{
				requestBody: `{
									"repaccesstype": "public"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom+1)},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_createRHandler ----------> create repository 4 ",
			param: param{
				requestBody: `{
									"comment": "中国移动北京终端详情"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom+2)},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_createRHandler ----------> create repository 5 ",
			param: param{
				requestBody: `{
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
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom+3)},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_createRHandler ----------> create repository 6 ",
			param: param{
				requestBody: `{
									"label": {
										"sys": {
											"loc": "北京"
											},
										"opt": {
											"age": 22
											}
										}
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom+4)},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_createRHandler ----------> create repository 7 ",
			param: param{
				requestBody: `{

								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom+5)},
				db:         db.copy(),
				login_name: USERNAME,
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
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("POST", "/repositories/rep0001", strings.NewReader(p.requestBody))
		get(err)
		code, msg := createRHandler(r, p.rsp, p.param, p.db, p.login_name)

		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}
		t.Log("")
	}
}

func Test_updateRHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: "Test_updateRHandler ----------> update repository 1 ",
			param: param{
				requestBody: `{
									"repaccesstype": "private"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom)},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    fmt.Sprintf("app0001_%d", ramdom),
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_updateRHandler ----------> update repository 2 ",
			param: param{
				requestBody: `{
									"comment": "中国移动北京终端详情2"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom)},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    fmt.Sprintf("app0001_%d", ramdom),
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_updateRHandler ----------> update repository 3 ",
			param: param{
				requestBody: `{
									"comment": "中国移动北京终端详情3"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": "repnameNotExist"},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    "repnameNotExist",
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).message, fmt.Sprintf(" %s=%s", COL_REPNAME, "repnameNotExist")),
				}},
			},
		},
		Context{
			description: "Test_updateRHandler ----------> update repository 4 ",
			param: param{
				requestBody: `{
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom)},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    fmt.Sprintf("app0001_%d", ramdom),
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
	t.Log("---------------------------------------------------------")
	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("GET", fmt.Sprintf("/repositories/%s", p.repName), strings.NewReader(p.requestBody))
		get(err)

		r.Header.Set("User", USERNAME)

		code, msg := updateRHandler(r, p.rsp, p.param, p.db, p.login_name)

		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}

	}
}

func Test_getRHandler(t *testing.T) {

	contexts := []Context{
		Context{
			description: "Test_getRHandler ----------> get repository 1 ",
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom)},
				db:          db.copy(),
				login_name:  USERNAME,
				repName:     fmt.Sprintf("app0001_%d", ramdom),
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_getRHandler ----------> get repository 2 ",
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": "repnameNotExist"},
				db:          db.copy(),
				login_name:  USERNAME,
				repName:     "repnameNotExist",
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).message, fmt.Sprintf(" %s=%s", COL_REPNAME, "repnameNotExist")),
				}},
			},
		},
	}
	t.Log("---------------------------------------------------------")
	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("GET", fmt.Sprintf("/repositories/%s", p.repName), strings.NewReader(p.requestBody))
		get(err)

		r.Header.Set("User", USERNAME)

		code, msg := getRHandler(r, p.rsp, p.param, p.db)

		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}

	}
}

func Test_delRHandler(t *testing.T) {

	contexts := []Context{
		Context{
			description: "Test_getRHandler ----------> delete repository 1 ",
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": fmt.Sprintf("app0001_%d", ramdom)},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    fmt.Sprintf("app0001_%d", ramdom),
			},
			expect: expect{
				code: 200,
				body: Body{Result{
					Code: 0,
					Msg:  "OK",
				}},
			},
		},
		Context{
			description: "Test_getRHandler ----------> delete repository 2 ",
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": "repnameNotExist"},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    "repnameNotExist",
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).message, fmt.Sprintf(" %s=%s", COL_REPNAME, "repnameNotExist")),
				}},
			},
		},
	}
	t.Log("---------------------------------------------------------")
	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("DELETE", fmt.Sprintf("/repositories/%s", p.repName), strings.NewReader(p.requestBody))
		get(err)

		r.Header.Set("User", USERNAME)

		code, msg := delRHandler(r, p.rsp, p.param, p.db, p.login_name)
		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}

	}
}

func (expect *expect) expect(t *testing.T, resutlCode int, resutlData string) bool {
	if expect.code != resutlCode {
		t.Errorf("expected http.code:%d != return http.code:%d", expect.code, resutlCode)
		return false
	}

	res := new(Result)
	json.Unmarshal([]byte(resutlData), res)

	if expect.body.Code != res.Code {
		t.Errorf("expected http.Code(%d) != return http.Code(%d)", expect.body.Code, res.Code)
		return false
	}

	if expect.body.Msg != res.Msg {
		t.Errorf("expected http.Msg(%s) != return http.Msg(%s)", expect.body.Msg, res.Msg)
		return false
	}

	//	if expect.body.Data != res.Data {
	//		t.Errorf("expected http.Data(%+v) != return http.Data(%+v)", expect.body.Data, res.Data)
	//		return false
	//	}

	return true
}
