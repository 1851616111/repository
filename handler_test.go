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

var (
	ramdom    int
	repnames  []string
	itemnames []string
)

func init() {
	rd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	ramdom = rd.Int()

	for i := 1; i <= 7; i++ {
		repnames = append(repnames, initRepositoryName(i))
		itemnames = append(itemnames, initDataitemName(i))
	}
	go q_c.serve(&db)
}

func Test_createRHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.创建Repository(全部参数) ----------> %s", repnames[0]),
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
			description: fmt.Sprintf("2.创建Repository(参数repository_name重复) ----------> %s", repnames[0]),
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
				login_name: USERNAME,
			},
			expect: expect{
				code: 400,
				body: Body{Result{
					Code: 1008,
					Msg:  errRepDatabaseOperate(repnames[0]),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.创建Repository(参数只有repaccesstype) ----------> %s", repnames[2]),
			param: param{
				requestBody: `{
									"repaccesstype": "public"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[2]},
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
			description: fmt.Sprintf("4.创建Repository(参数只有comment) ----------> %s", repnames[3]),
			param: param{
				requestBody: `{
									"comment": "中国移动北京终端详情"
								}`,
				rsp:        &Rsp{w: httptest.NewRecorder()},
				param:      martini.Params{"repname": repnames[3]},
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
			description: fmt.Sprintf("5.创建Repository(参数只有label) ----------> %s", repnames[4]),
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
			description: fmt.Sprintf("6.创建Repository(参数只有label) ----------> %s", repnames[5]),
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
			description: fmt.Sprintf("7.创建Repository(参数为空) ----------> %s", repnames[6]),
			param: param{
				requestBody: `{}`,
				rsp:         &Rsp{w: httptest.NewRecorder()},
				param:       martini.Params{"repname": repnames[6]},
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

func Test_createDHandler(t *testing.T) {
	contexts := []Context{
		Context{
			description: fmt.Sprintf("1.创建Dataitem(全部参数) ----------> %s/%s", repnames[0], itemnames[0]),
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
			description: fmt.Sprintf("2.创建Dataitem(重复dataitem) ----------> %s/%s", repnames[0], itemnames[0]),
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
					Msg:  errItemDatabaseOperate(repnames[0], itemnames[0]),
				}},
			},
		},
		Context{
			description: fmt.Sprintf("3.创建Dataitem(不存在的repository) ----------> %s/%s", repnames[1], itemnames[1]),
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
			description: fmt.Sprintf("4.创建Dataitem(必选参数label.sys.supply_style缺失) ----------> %s/%s", repnames[2], itemnames[2]),
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
			description: fmt.Sprintf("5.创建Dataitem(必选参数label.sys.supply_style违法) ----------> %s/%s", repnames[2], itemnames[2]),
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
			description: fmt.Sprintf("6.创建Dataitem(label自定义参数) ----------> %s/%s", repnames[2], itemnames[2]),
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
			description: fmt.Sprintf("6.创建Dataitem(meta,sample,comment不传) ----------> %s/%s", repnames[2], itemnames[3]),
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
		code, msg := delDHandler(r, p.rsp, p.param, p.db, p.login_name)

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

	return true
}

func errRepDatabaseOperate(repositoryName string) string {
	return fmt.Sprintf("database operate : insertDocument :: caused by :: 11000 E11000 duplicate key error index: datahub.repository.$repository_name_1 dup key: { : \"%s\" }", repositoryName)
}
func errItemDatabaseOperate(repositoryName, itemName string) string {
	return fmt.Sprintf("database operate : insertDocument :: caused by :: 11000 E11000 duplicate key error index: datahub.dataitem.$repository_name_1_dataitem_name_1 dup key: { : \"%s\", : \"%s\" }", repositoryName, itemName)
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

func initRepositoryName(casenum int) string {
	return fmt.Sprintf("test_repository_%d_case_%d", ramdom, casenum)
}

func initDataitemName(casenum int) string {
	return fmt.Sprintf("test_dataitem_%d_case_%d", ramdom, casenum)
}
