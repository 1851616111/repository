package main

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub_repository/log"
	"github.com/go-martini/martini"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	USERNAME      = "panxy3@asiainfo.com"
	ADMINUSERNAME = USERNAME
)

var (
	ramdom         int
	repnames       []string
	itemnames      []string
	tagnames       []string
	selectlabel    string
	newselectlabel string
	token          string
)

func init() {
	Log = log.NewLogger("test")
	rd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	ramdom = rd.Int()

	for i := 1; i <= 7; i++ {
		repnames = append(repnames, initRepositoryName(i))
		itemnames = append(itemnames, initDataitemName(i))
		tagnames = append(tagnames, initTagName(i))
	}

	selectlabel = fmt.Sprintf("精选栏目_%d", ramdom)
	newselectlabel = fmt.Sprintf("精选栏目_new_%d", ramdom)

	db := initDB()

	go q_c.serve(&db)

	token = getToken(USERNAME, "q")
	if len(token) != 32 {
		Log.Error("Init token failed\n")
	} else {
		Log.Infof("Init token:%s success\n", token)
	}
}

func Test_createRHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【拥有者】【新增】Rep (全部参数) ----------> (%s)", repnames[0]),
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
				param:      martini.Params{"repname": repnames[0]},
				db:         db.copy(),
				limit:      Limit{Rep_Private: 100, Rep_Public: 100},
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
			description: fmt.Sprintf("2.【拥有者】【新增】Rep (参数repository_name重复) ----------> (%s)", repnames[0]),
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
				param:      martini.Params{"repname": repnames[0]},
				db:         db.copy(),
				limit:      Limit{Rep_Private: 100, Rep_Public: 100},
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1008,
					Msg:  duplicatedRep(repnames[0]),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.【拥有者】【新增】Rep (参数只有repaccesstype) ----------> (%s)", repnames[2]),
			param: param{
				requestBody: `{
									"repaccesstype": "public"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2]},
				db:         db.copy(),
				limit:      Limit{Rep_Private: 100, Rep_Public: 100},
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
			description: fmt.Sprintf("4.【拥有者】【新增】Rep (参数只有comment) ----------> (%s)", repnames[3]),
			param: param{
				requestBody: `{
									"comment": "中国移动北京终端详情"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[3]},
				db:         db.copy(),
				limit:      Limit{Rep_Private: 100, Rep_Public: 100},
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
			description: fmt.Sprintf("5.【拥有者】【新增】Rep (参数只有label) ----------> (%s)", repnames[4]),
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
				param:      martini.Params{"repname": repnames[4]},
				db:         db.copy(),
				limit:      Limit{Rep_Private: 100, Rep_Public: 100},
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
			description: fmt.Sprintf("6.【拥有者】【新增】Rep (参数只有label) ----------> (%s)", repnames[5]),
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
				param:      martini.Params{"repname": repnames[5]},
				db:         db.copy(),
				limit:      Limit{Rep_Private: 100, Rep_Public: 100},
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
			description: fmt.Sprintf("7.【拥有者】【新增】Rep (参数为空) ----------> (%s)", repnames[6]),
			param: param{
				requestBody: `{}`,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[6]},
				db:          db.copy(),
				limit:       Limit{Rep_Private: 100, Rep_Public: 100},
				login_name:  USERNAME,
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
		r.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
		r.Header.Set("User", p.login_name)
		code, msg := createRHandler(r, p.rsp, p.param, p.db, p.login_name, p.limit)

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

func Test_createDHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【拥有者】【新增】Item (全部参数) ----------> (%s/%s)", repnames[0], itemnames[0]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"meta": "样例数据",
								"sample": "元数据",
								"comment": "对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {
										"supply_style": "batch"
									},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
			description: fmt.Sprintf("2.【拥有者】【新增】Item (重复dataitem) ----------> (%s/%s)", repnames[0], itemnames[0]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"meta": "样例数据",
								"sample": "元数据",
								"comment": "对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {
										"supply_style": "batch"
									},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1008,
					Msg:  duplicatedItem(repnames[0], itemnames[0]),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.【拥有者】【新增】Item (不存在的repository) ----------> (%s/%s)", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"meta": "样例数据",
								"sample": "元数据",
								"comment": "对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {
										"supply_style": "batch"
									},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[1], "itemname": itemnames[1]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf("repname : %s", repnames[1])),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("4.【拥有者】【新增】Item (必选参数label.sys.supply_style缺失) ----------> (%s/%s)", repnames[2], itemnames[2]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"meta": "样例数据",
								"sample": "元数据",
								"comment": "对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2], "itemname": itemnames[2]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1400,
					Msg:  fmt.Sprintf("%s: %s", E(ErrorCodeNoParameter).Message, "label.sys.supply_style"),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("5.【拥有者】【新增】Item (必选参数label.sys.supply_style违法) ----------> (%s/%s)", repnames[2], itemnames[2]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"meta": "样例数据",
								"sample": "元数据",
								"comment": "对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {"supply_style": "batchinvalid"},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2], "itemname": itemnames[2]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1007,
					Msg:  fmt.Sprintf("%s: %s", E(ErrorCodeInvalidParameters).Message, "label.sys.supply_style"),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("6.【拥有者】【新增】Item (label自定义参数) ----------> (%s/%s)", repnames[2], itemnames[2]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"meta": "样例数据",
								"sample": "元数据",
								"comment": "对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {
											"supply_style": "batch"
											},
									"opt": {},
									"owner": {
												"key":true,
												"param":null
											},
									"other": {
												"key":true,
												"param":null
											}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2], "itemname": itemnames[2]},
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
			description: fmt.Sprintf("7.【拥有者】【新增】Item (meta,sample,comment不传) ----------> (%s/%s)", repnames[2], itemnames[3]),
			param: param{
				requestBody: `{
								"itemaccesstype": "private",
								"label": {
									"sys": {
											"supply_style": "batch"
											},
									"opt": {},
									"owner": {
												"key":true,
												"param":null
											},
									"other": {
												"key":true,
												"param":null
											}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2], "itemname": itemnames[3]},
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
		r, err := http.NewRequest("POST", "/repositories/rep/item", strings.NewReader(p.requestBody))
		get(err)
		code, msg := createDHandler(r, p.rsp, p.param, p.db, p.login_name)

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

func Test_searchHandler(t *testing.T) {
	contexts := []Context{
		Context{

			description: fmt.Sprintln("1.【任意】【查询】Search 关键字 (mobile)"),
			param: param{
				rsp: &Rsp{w: httptest.NewRecorder()},
				db:  db.copy(),
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
		r, err := http.NewRequest("GET", fmt.Sprintf("/search?text=mobile", p.repName), strings.NewReader(p.requestBody))
		get(err)

		code, msg := searchHandler(r, p.rsp, p.db)

		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}
	}
}

func Test_createTagHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【拥有者】【新增】Tag (全部参数) ----------> (%s/%s/%s)", repnames[0], itemnames[0], tagnames[0]),
			param: param{
				requestBody: `{
								"comment":"2001MB"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[0]},
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
			description: fmt.Sprintf("2.【拥有者】【新增】Tag (重复创建) ----------> (%s/%s/%s)", repnames[0], itemnames[0], tagnames[0]),
			param: param{
				requestBody: `{
								"comment":"20001MB"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1008,
					Msg:  duplicatedTag(repnames[0], itemnames[0], tagnames[0]),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.【拥有者】【新增】Tag (所属dataitem不存在) ----------> (%s/%s/%s)", repnames[0], itemnames[2], tagnames[1]),
			param: param{
				requestBody: `{
								"comment":"20001MB"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[2], "tag": tagnames[1]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf("itemname : %s", itemnames[2])),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("4.【拥有者】【新增】Tag (请求body未传参数) ----------> (%s/%s/%s)", repnames[0], itemnames[0], tagnames[1]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[1]},
				db:          db.copy(),
				login_name:  USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1400,
					Msg:  fmt.Sprintf("%s: %s", E(ErrorCodeNoParameter).Message, ""),
				}},
			},
		},
	}

	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("POST", "/repositories/rep/item/tag", strings.NewReader(p.requestBody))
		get(err)
		code, msg := createTagHandler(r, p.rsp, p.param, p.db, p.login_name, nil)

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

func Test_setSelectLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【管理员】【新增】SelectLabel --------------> (%s)", selectlabel),
			param: param{
				requestBody: `{
							"order": 1,
							"icon":"path1"
						}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"labelname": selectlabel},
				db:         db.copy(),
				login_name: ADMINUSERNAME,
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
		r, err := http.NewRequest("POST", "/select_labels/:labelname", strings.NewReader(p.requestBody))
		get(err)
		code, msg := setSelectLabelHandler(r, p.rsp, p.param, p.db)

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

func Test_getSelectLabelsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【管理员】【查询】SelectLabel --------------> (%s)", selectlabel),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"labelname": selectlabel},
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
		r, err := http.NewRequest("GET", "/select_labels", strings.NewReader(p.requestBody))
		get(err)
		code, msg := getSelectLabelsHandler(r, p.rsp, p.db)

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

func Test_updateSelectLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【管理员】【更新】SelectLabel   --------------> (old: %s, new : %s)", selectlabel, newselectlabel),
			param: param{
				requestBody: fmt.Sprintf(`{
							"order": 2,
							"icon":"path2",
							"newlabelname":"%s"
						}`, newselectlabel),
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"labelname": selectlabel},
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
		r, err := http.NewRequest("PUT", "/select_labels/labelname", strings.NewReader(p.requestBody))
		get(err)
		code, msg := updateSelectLabelHandler(r, p.rsp, p.param, p.db)

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

func Test_updateSelectHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【管理员】将Item添加至精选 ----------> (%s/%s)", repnames[0], itemnames[0]),
			param: param{

				requestForm: url.Values{"select_labels": []string{newselectlabel}, "order": []string{"100"}},
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("POST", "/selects/repname/itemname", strings.NewReader(p.requestBody))
		get(err)
		r.PostForm = p.requestForm
		code, msg := updateSelectHandler(r, p.rsp, p.param, p.db)

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

func Test_getSelectsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【任意】【查询】精选的Item(按照order排序)"),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
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
		r, err := http.NewRequest("GET", "/selects?select_labels=selectlabel1", strings.NewReader(p.requestBody))
		get(err)
		code, msg := getSelectsHandler(r, p.rsp, p.db)

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

func Test_delSelectHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【管理员】删除精选内容"),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				db:         db.copy(),
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				login_name: ADMINUSERNAME,
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
		r, err := http.NewRequest("DELETE", "/selects/rep1/dataitem1", strings.NewReader(p.requestBody))
		get(err)
		code, msg := getSelectsHandler(r, p.rsp, p.db)

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

func Test_upsertRLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.新增/更新repository的label的某条属性 ----------> %s ", repnames[0]),
			param: param{
				requestForm: url.Values{"owner.age": []string{"15"}},
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("PUT", "", strings.NewReader(p.requestBody))
		get(err)
		r.Header.Set("User", p.login_name)
		code, msg := upsertRLabelHandler(r, p.rsp, p.param, p.db)

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

func Test_upsertDLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.新增/更新dataitem的label的某条属性  ----------> %s/%s/", repnames[0], itemnames[0]),
			param: param{
				requestForm: url.Values{"owner.age": []string{"16"}},
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("PUT", "", strings.NewReader(p.requestBody))
		get(err)
		r.Header.Set("User", p.login_name)
		code, msg := upsertDLabelHandler(r, p.rsp, p.param, p.db)

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

func Test_delRLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.删除repository的label的某条属性 ----------> %s ", repnames[0]),
			param: param{
				requestForm: url.Values{"owner.age": []string{"15"}},
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("DELETE", "", strings.NewReader(p.requestBody))
		get(err)
		r.Header.Set("User", p.login_name)
		code, msg := delRLabelHandler(r, p.rsp, p.param, p.db)

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

func Test_delDLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.删除dataitem的label的某条属性  ----------> %s/%s", repnames[0], itemnames[0]),
			param: param{
				requestForm: url.Values{"owner.age": []string{"16"}},
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("DELETE", "", strings.NewReader(p.requestBody))
		get(err)
		r.Header.Set("User", p.login_name)
		code, msg := delDLabelHandler(r, p.rsp, p.param, p.db)

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

func Test_setRepPmsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【rep拥有者】将某用户（非自己）加入或更新rep白名单中 ----------> %s", repnames[0]),
			param: param{
				requestBody: `{"username":"chai@asiainfo.com","opt_permission":1}`,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("POST", "/permission/rep", strings.NewReader(p.requestBody))
		get(err)
		Rep_Permission := Rep_Permission{Repository_name: repnames[0]}
		code, msg := setRepPmsHandler(r, p.rsp, p.param, p.db, Rep_Permission)

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

func Test_setItemPmsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【item拥有者】将某用户（非自己）加入或更新item白名单中 ----------> %s", repnames[0]),
			param: param{
				requestBody: `{"username":"chai@asiainfo.com"}`,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("POST", "/permission/rep/item", strings.NewReader(p.requestBody))
		get(err)
		item_permission := Item_Permission{Repository_name: repnames[0], Dataitem_name: itemnames[0]}
		code, msg := setItemPmsHandler(r, p.rsp, p.param, p.db, item_permission)

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

func Test_getRepPmsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【rep拥有者】查询自己rep白名单的username列表，及相应的读写权限 ----------> %s", repnames[0]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0]},
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
		r, err := http.NewRequest("GET", "/permission/rep", strings.NewReader(p.requestBody))
		get(err)
		Rep_Permission := Rep_Permission{Repository_name: repnames[0]}
		code, msg := getRepPmsHandler(r, p.rsp, p.param, p.db, Rep_Permission)

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

func Test_getItemPmsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【item拥有者】将某用户（非自己）加入或更新item白名单中 ----------> %s", repnames[0], itemnames[0]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
		r, err := http.NewRequest("GET", "/permission/rep/item", strings.NewReader(p.requestBody))
		get(err)
		item_permission := Item_Permission{Repository_name: repnames[0], Dataitem_name: itemnames[0]}
		code, msg := getItemPmsHandler(r, p.rsp, p.param, p.db, item_permission)

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

func Test_delRepPmsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【rep拥有者】将某用户（非自己）从rep白名单中删除 ----------> %s", repnames[0]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0]},
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
		r, err := http.NewRequest("DELETE", fmt.Sprintf("/permission/%s?username=chai@asiainfo.com", repnames[0]), strings.NewReader(p.requestBody))
		get(err)
		Rep_Permission := Rep_Permission{Repository_name: repnames[0]}
		code, msg := delRepPmsHandler(r, p.rsp, p.param, p.db, Rep_Permission)

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

func Test_delItemPmsHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【item拥有者】将某用户（非自己）从item白名单中移除 ----------> %s", repnames[0], itemnames[0]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
		r, err := http.NewRequest("DELETE", fmt.Sprintf("/permission/%s/%s?username=chai@asiainfo.com", repnames[0], itemnames[0]), strings.NewReader(p.requestBody))
		get(err)

		item_permission := Item_Permission{Repository_name: repnames[0], Dataitem_name: itemnames[0]}
		code, msg := delItemPmsHandler(r, p.rsp, p.param, p.db, item_permission)

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
			description: fmt.Sprintf("1.更新Repository(repaccesstype) ----------> %s", repnames[0]),
			param: param{
				requestBody: `{
									"repaccesstype": "private"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[0],
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
			description: fmt.Sprintf("2.更新Repository(comment) ----------> %s", repnames[0]),
			param: param{
				requestBody: `{
									"comment": "中国移动北京终端详情2"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[0],
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
			description: fmt.Sprintf("3.更新Repository(不存在repository) ----------> %s", "repnameNotExist"),
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
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf(" %s=%s", COL_REPNAME, "repnameNotExist")),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("4.更新Repository(参数为空) ----------> %s", repnames[0]),
			param: param{
				requestBody: `{
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[0],
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

func Test_updateDHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.更新Dataitem(全部参数) ----------> %s/%s", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"itemaccesstype": "public",
								"meta": "样例数据更新",
								"sample": "元数据更新",
								"comment": "更新对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {
										"supply_style": "api"
									},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
			description: fmt.Sprintf("1.更新Dataitem(全部参数) ----------> %s/%s", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"itemaccesstype": "public",
								"meta": "样例数据更新",
								"sample": "元数据更新",
								"comment": "更新对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。",
								"label": {
									"sys": {
										"supply_style": "api"
									},
									"opt": {},
									"owner": {},
									"other": {}
								}
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
			description: fmt.Sprintf("2.更新Dataitem(更新其中一个参数itemaccesstype) ----------> %s/%s", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"itemaccesstype": "public"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
			description: fmt.Sprintf("3.更新Dataitem(更新其中一个参数meta) ----------> %s/%s", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"meta": "样例数据更新"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
			description: fmt.Sprintf("4.更新Dataitem(更新其中一个参数sample) ----------> %s/%s", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"sample": "元数据更新"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
			description: fmt.Sprintf("5.更新Dataitem(更新其中一个参数comment) ----------> %s/%s", repnames[1], itemnames[1]),
			param: param{
				requestBody: `{
								"comment": "更新对终端使用情况、变化情况进行了全方面的分析。包括分品牌统计市场存量、新增、机型、数量、换机等情况。终端与ARPU、DOU、网龄的映射关系。终端的APP安装情况等。"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
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
		r, err := http.NewRequest("PUT", "/repositories/rep/item", strings.NewReader(p.requestBody))
		get(err)
		code, msg := updateDHandler(r, p.rsp, p.param, p.db, p.login_name)

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

func Test_updateTagHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.更新Tag(全部参数) ----------> %s/%s/%s", repnames[0], itemnames[0], tagnames[0]),
			param: param{
				requestBody: `{
								"comment":"update2001MB"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[0]},
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
			description: fmt.Sprintf("2.更新Tag(所属dataitem不存在) ----------> %s/%s/%s", repnames[0], itemnames[1], tagnames[0]),
			param: param{
				requestBody: `{
								"comment":"20001MB"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[1], "tag": tagnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf("itemname : %s", itemnames[1])),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.更新Tag(更新的tag不存在) ----------> %s/%s/%s", repnames[0], itemnames[0], tagnames[2]),
			param: param{
				requestBody: `{
								"comment":"20001MB"
							}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[2]+"123"},
				db:         db.copy(),
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf("tagname : %s", tagnames[2])),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("4.更新Tag(请求body未传参数) ----------> %s/%s/%s", repnames[0], itemnames[0], tagnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1400,
					Msg:  fmt.Sprintf("%s: %s", E(ErrorCodeNoParameter).Message, ""),
				}},
			},
		},
	}

	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("PUT", "/repositories/rep/item/tag", strings.NewReader(p.requestBody))
		get(err)
		code, msg := updateTagHandler(r, p.rsp, p.param, p.db, p.login_name)

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

func Test_getRHandler(t *testing.T) {
	contexts := []Context{
		Context{

			description: fmt.Sprintf("1.查询Repository ----------> %s", repnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
				repName:     repnames[0],
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
			description: fmt.Sprintf("2.查询Repository(不存在repository) ----------> %s", repnames[0]),
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
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf(" %s=%s", COL_REPNAME, "repnameNotExist")),
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

func Test_getDHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.查询Dataitem ----------> %s/%s", repnames[0], itemnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				db:          db.copy(),
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
			description: fmt.Sprintf("2.查询Dataitem(所在repository不存在) ----------> %s/%s", repnames[1], itemnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[1], "itemname": itemnames[0]},
				db:          db.copy(),
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf(" %s=%s", COL_REPNAME, repnames[1])),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.查询Dataitem(查询Dataitem不存在) ----------> %s/%s", repnames[0], itemnames[5]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[5]},
				db:          db.copy(),
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf(" %s=%s,%s=%s ", COL_REPNAME, repnames[0], COL_ITEM_NAME, itemnames[5])),
				}},
			},
		},
	}

	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("GET", "/repositories/rep/item", strings.NewReader(p.requestBody))
		get(err)
		r.Header.Set("User", USERNAME)
		code, msg := getDHandler(r, p.rsp, p.param, p.db)

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

func Test_delTagHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.删除Tag(已存在) ----------> %s/%s/%s", repnames[0], itemnames[0], tagnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
			description: fmt.Sprintf("2.删除Tag(所属dataitem不存在) ----------> %s/%s/%s", repnames[0], itemnames[1], tagnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[1], "tag": tagnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf("itemname : %s", itemnames[1])),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.删除Tag(dataitem存在,tag不存在) ----------> %s/%s/%s", repnames[0], itemnames[0], tagnames[1]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0], "tag": tagnames[1]},
				db:          db.copy(),
				login_name:  USERNAME,
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
		r, err := http.NewRequest("DELETE", "/repositories/rep/item/tag", strings.NewReader(p.requestBody))
		get(err)
		r.Header.Set("User", USERNAME)
		code, msg := delTagHandler(r, p.rsp, p.param, p.db, p.login_name, nil)

		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}
		t.Log("")
	}

	time.Sleep(time.Second)
}

func Test_delDHandler(t *testing.T) {

	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.删除Dataitem(已存在dataitem) ----------> %s/%s", repnames[0], itemnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[0], "itemname": itemnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
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
			description: fmt.Sprintf("2.删除Dataitem(已存在dataitem) ----------> %s/%s", repnames[2], itemnames[2]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[2], "itemname": itemnames[2]},
				db:          db.copy(),
				login_name:  USERNAME,
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
			description: fmt.Sprintf("3.删除Dataitem(已存在dataitem) ----------> %s/%s", repnames[2], itemnames[3]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[2], "itemname": itemnames[3]},
				db:          db.copy(),
				login_name:  USERNAME,
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
			description: fmt.Sprintf("4.删除Dataitem(不存在的dataitem) ----------> %s/%s", repnames[1], itemnames[0]),
			param: param{
				requestBody: ``,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[1], "itemname": itemnames[0]},
				db:          db.copy(),
				login_name:  USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1009,
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf(" %s=%s %s:=%s", COL_REPNAME, repnames[1], COL_ITEM_NAME, itemnames[0])),
				}},
			},
		},
	}

	for _, v := range contexts {
		p := v.param
		expect := v.expect
		r, err := http.NewRequest("DELETE", "/repositories/rep/item", strings.NewReader(p.requestBody))
		get(err)
		code, msg := delDHandler(r, p.rsp, p.param, p.db, p.login_name, nil)

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

func Test_delRHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.删除Repository ----------> %s", repnames[0]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[0]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[0],
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
			description: fmt.Sprintf("2.删除Repository ----------> %s", "repnameNotExist"),
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
					Msg:  fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, fmt.Sprintf(" %s=%s", COL_REPNAME, "repnameNotExist")),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("2.删除Repository ----------> %s", repnames[2]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[2],
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
			description: fmt.Sprintf("3.删除Repository ----------> %s", repnames[3]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[3]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[3],
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
			description: fmt.Sprintf("4.删除Repository ----------> %s", repnames[4]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[4]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[4],
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
			description: fmt.Sprintf("5.删除Repository ----------> %s", repnames[5]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[5]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[5],
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
			description: fmt.Sprintf("6.删除Repository ----------> %s", repnames[6]),
			param: param{
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[6]},
				db:         db.copy(),
				login_name: USERNAME,
				repName:    repnames[6],
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
		r, err := http.NewRequest("DELETE", fmt.Sprintf("/repositories/%s", p.repName), strings.NewReader(p.requestBody))
		get(err)

		r.Header.Set("User", USERNAME)

		code, msg := delRHandler(r, p.rsp, p.param, p.db, p.login_name, nil)
		if !expect.expect(t, code, msg) {
			t.Logf("%s fail.", v.description)
			t.Log(code)
			t.Log(msg)
		} else {
			t.Logf("%s success.", v.description)
		}

	}
}

func Test_delSelectLabelHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.【管理员】删除精选栏目--------------> %s", selectlabel),
			param: param{
				rsp:   &Rsp{w: httptest.NewRecorder()},
				param: martini.Params{"labelname": newselectlabel},
				db:    db.copy(),
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
		r, err := http.NewRequest("DELETE", "/select_labels/:labelname", strings.NewReader(p.requestBody))
		get(err)
		code, msg := delSelectLabelHandler(r, p.rsp, p.param, p.db.copy())

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

	return true
}

func duplicatedRep(repositoryName string) string {
	return fmt.Sprintf("database operate : insertDocument :: caused by :: 11000 E11000 duplicate key error index: datahub.repository.$repository_name_1  dup key: { : \"%s\" }", repositoryName)
}

func duplicatedItem(repositoryName, itemName string) string {
	return fmt.Sprintf("database operate : insertDocument :: caused by :: 11000 E11000 duplicate key error index: datahub.dataitem.$repository_name_1_dataitem_name_1  dup key: { : \"%s\", : \"%s\" }", repositoryName, itemName)
}

func duplicatedTag(repositoryName, itemName, tagName string) string {
	return fmt.Sprintf("database operate : insertDocument :: caused by :: 11000 E11000 duplicate key error index: datahub.tag.$repository_name_1_dataitem_name_1_tag_1  dup key: { : \"%s\", : \"%s\", : \"%s\" }", repositoryName, itemName, tagName)
}

type param struct {
	requestForm url.Values
	requestBody string
	rsp         *Rsp
	param       martini.Params
	db          *DB
	login_name  string
	limit       Limit
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

func initRepositoryName(casenum int) string {
	return fmt.Sprintf("test_repository_%d_case_%d", ramdom, casenum)
}

func initDataitemName(casenum int) string {
	return fmt.Sprintf("test_dataitem_%d_case_%d", ramdom, casenum)
}

func initTagName(casenum int) string {
	return fmt.Sprintf("test_tag_%d_case_%d", ramdom, casenum)[0:20]
}
