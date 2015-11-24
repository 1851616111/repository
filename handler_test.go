package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"strings"
	"testing"
)

func init() {
	se := connect(DB_URL_MONGO)
	db = DB{*se}
}

type request struct {
	requestBody string
	rsp         *Rsp
	param       martini.Params
	db          *DB
	login_name  string
}

func (res *Result) expect(t testing.T, target *Result) bool {
	if res.Code != target.Code {
		t.Errorf("expected resutl.code:%d != return resutl.code:%d", res.Code, target.Code)
		return false
	}
	if res.Msg != target.Msg {
		t.Errorf("expected resutl.msg :%s != return resutl.msg:%s", res.Msg, target.Msg)
		return false
	}
	if res.Data != target.Data {
		t.Errorf("expected resutl.data:%+v != return resutl.data:%+v", res.Data, target.Data)
		return false
	}
	return true
}

func TestcreateRHandler(t *testing.T) {

	contexts := []struct {
		req      request
		expected Result
	}{
		{
			req: request{
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
				rsp:        nil,
				param:      martini.Params{"repname": "app0001"},
				db:         db,
				login_name: "panxy3@asiainfo.com",
			},

			expected: Result{Code: 200, Msg: "OK"},
		},
	}

	for _, v := range contexts {
		req := v.req
		r, err := http.NewRequest("POST", "/repositories/rep0001", strings.NewReader(req.requestBody))
		get(err)
		code, msg := createRHandler(r, req.rsp, req.param, req.db, req.login_name)
		res := Result{Code: code, Msg: msg}
		v.expected.expect(t, res)
	}
}
