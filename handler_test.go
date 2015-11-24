package main

import (
	"github.com/go-martini/martini"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type testParam struct {
	requestBody string
	rsp         *Rsp
	param       martini.Params
	db          *DB
	login_name  string
}

func init() {
	//	rd := rand.New(rand.NewSource(time.Now.Nanosecond()))

}

func Test_createRHandler(t *testing.T) {

	contexts := []struct {
		description string
		param       testParam
		expected    Result
	}{
		{
			description: "----------> create repository",
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
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": "app0001"},
				db:         &db,
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
		log.Println(code)
		log.Println(msg)

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
