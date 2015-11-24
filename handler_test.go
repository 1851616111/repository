package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"strings"
	"testing"
	"net/http/httptest"
)

//var (
//	database DB
//)
//
//func init() {
//
//}

type testParam struct {
	requestBody string
	rsp         *Rsp
	param       martini.Params
	db          *DB
	login_name  string
}

func Test_createRHandler(t *testing.T) {
	se := connect(fmt.Sprintf(`%s:%s/datahub?maxPoolSize=50`, "10.1.235.98", "27017"))
	database := DB{*se}

	contexts := []struct {
		param testParam
		expected Result
	}{
		{
			param: testParam{
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
				rsp:        &Rsp{w:httptest.NewRecorder()},
				param:      martini.Params{"repname": "app0001"},
				db:         &database,
				login_name: "panxy3@asiainfo.com",
			},

			expected: Result{Code: 200, Msg: "OK"},
		},
	}

	for _, v := range contexts {
		p := v.param
		r, err := http.NewRequest("POST", "/repositories/rep0001", strings.NewReader(p.requestBody))
		get(err)
		code, msg := createRHandler(r, p.rsp, p.param, p.db, p.login_name)
		//		res := Result{Code: code, Msg: msg}
		//		v.expected.expect(t, res)

		t.Log(code)
		t.Log(msg)

	}
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
