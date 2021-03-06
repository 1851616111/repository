package main

import (
	"testing"
)

func Test_buildTime(t *testing.T) {
	param := []string{"2013-01-02 15:04:05", "2014-01-02 15:04:05", "2015-05-26 15:04:05", "2015-10-28 15:04:05", "2015-11-28 15:04:05", "2015-11-30 15:04:05"}
	for i := 0; i < len(param); i++ {
		t.Log(param[i])
		t.Log(buildTime(param[i]))
	}
}

func Test_ifInLabel(t *testing.T) {
	label := map[string]interface{}{"sys": map[string]interface{}{"supply_style": "api"}}
	if err := ifInLabel(label, "supply_style"); err != nil {
		t.Errorf("Test_ifInLabel fail: %s", err.ErrToString())
	}
}

func Test_chkPrice(t *testing.T) {
	context := M{
		"api": []map[string]interface{}{
			map[string]interface{}{"units": 1000, "money": 5, "expire": DATAITEM_PRICE_EXPIRE},
			map[string]interface{}{"units": 10000, "money": 45, "expire": DATAITEM_PRICE_EXPIRE},
			map[string]interface{}{"units": 100000, "money": 400.00, "expire": DATAITEM_PRICE_EXPIRE},
		},
	}

	for _, v := range context {
		if err := chkPrice(v); err != nil {
			t.Errorf("Test_chkPrice fail: %s", err.ErrToString())
		}
	}
}

//func Test_getMd5(t *testing.T) {
//	in := []string{"88888888"}
//	expect := []string{"8ddcff3a80f4189ca1c9d4d902c3c909"}
//
//	for i, v := range in {
//		out := getMd5(v)
//		if out != expect[i] {
//			t.Errorf("input %s, output %s, expect %s", in[i], out, expect[i])
//		}
//	}
//}
//
//func Test_base64Encode(t *testing.T) {
//	in := []string{"panxy3@asiainfo.com:8ddcff3a80f4189ca1c9d4d902c3c909"}
//	expect := []string{"cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ=="}
//
//	for i, v := range in {
//		out := string(base64Encode([]byte(v)))
//		if out != expect[i] {
//			t.Errorf("Input: %s\n Output %s\n Expect %s\n", in[i], out, expect[i])
//		}
//	}
//}
//
//func Test_getToken(t *testing.T) {
//	user := "panxy3@asiainfo.com"
//	passwd := "q"
//	result := "a189775949e417acd7d4349de8e33000"
//
//	out := getToken(user, passwd)
//	if len(out) != len(result) {
//		t.Errorf("Input: %s\n Output %s\n Expect %s\n", user, passwd, result)
//	}
//}
