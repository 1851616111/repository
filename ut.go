package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TimeFormat             = "2006-01-02 15:04:05"
	TimeFormatDay          = "2006-01-02"
	DATAITEM_PRICE_EXPIRE  = 30
	DATAITEM_PRICE_MAX     = 6
	DEFINE_TAG_NAME        = "column"
	PRICE_STATE_FREE       = "免费"
	PRICE_STATE_FREE_LIMIT = "限量免费"
	PRICE_STATE_NOT_FREE   = "付费"
)

var (
	LOCAL_LOCATION     *time.Location
	COL_LABEL_CHILDREN = []string{"sys", "opt", "owner", "other"}
)

func init() {
	loc, err := time.LoadLocation("Local")
	chk(err)
	LOCAL_LOCATION = loc
}

type MM map[interface{}]M
type M map[interface{}]interface{}
type Ms map[string]interface{}
type Arr []interface{}

type Q struct {
	Columns    []string
	Conditions M
}

type Dim struct {
	sync.RWMutex
	mm MM
}

func (p *Dim) GetM(field string) (m M, exists bool) {
	p.RLock()
	defer p.RUnlock()
	m, exists = p.mm[field]
	return
}

func (p *Dim) Set(field, id, name string) {
	p.Lock()
	defer p.Unlock()
	p.mm[field][id] = name
}

func (p *Dim) SetM(field string, m M) {
	p.Lock()
	defer p.Unlock()
	p.mm[field] = m
}

type Rsp struct {
	w http.ResponseWriter
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
func get(err error) {
	if err != nil {
		Log.Error(err)
	}
}

type Result struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (p *Rsp) Json(code int, e *Error, data ...interface{}) (int, string) {
	//	p.w.Header().Set("Access-Control-Allow-Origin", "*")
	//	p.w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	//	p.w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept,X-Requested-With")
	p.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	p.w.Header().Set("Transfer", "-1")
	//	p.w.Header().Set("Transfer-Encoding", "chunked")

	result := new(Result)
	if len(data) > 0 {
		result.Data = data[0]
	}
	result.Code = e.Code
	result.Msg = e.Message

	b, err := json.Marshal(result)
	chk(err)
	return code, string(b)
}
func JsonResult(w http.ResponseWriter, statusCode int, e *Error, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if e == nil {
		e = ErrorNone
	}

	//result := Result {code: e.code, msg: e.message, data: data}
	result := Result{Code: e.Code, Msg: e.Message}

	jsondata, err := json.Marshal(&result)
	if err != nil {
		w.Write([]byte(getJsonBuildingErrorJson()))
	} else {
		w.Write(jsondata)
	}
}

var Json_ErrorBuildingJson []byte

func getJsonBuildingErrorJson() []byte {
	if Json_ErrorBuildingJson == nil {
		Json_ErrorBuildingJson = []byte(fmt.Sprintf(`{"code": %d, "msg": %s}`, ErrorJsonBuilding.Code, ErrorJsonBuilding.Message))
	}

	return Json_ErrorBuildingJson
}

func (p *Select) ParseRequeset(r *http.Request) error {
	t := reflect.TypeOf(*p)
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name := r.PostFormValue(strings.ToLower(f.Name)); name != "" {
			switch f.Type.Name() {
			case "int":
				i, _ := strconv.Atoi(name)
				v.FieldByName(f.Name).SetInt(int64(i))
			case "string":
				v.FieldByName(f.Name).SetString(name)
			case "float32":
				fallthrough
			case "float64":
				ff, _ := strconv.ParseFloat(name, 10)
				v.FieldByName(f.Name).SetFloat(ff)
			}
		} else {
			return errors.New(fmt.Sprintf("parse request err, no param: %s", strings.ToLower(f.Name)))
		}
	}
	return nil
}

func (p *repository) ParseRequeset(r *http.Request) error {
	t := reflect.TypeOf(*p)
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name := r.PostFormValue(strings.ToLower(f.Name)); name != "" {
			switch f.Type.Name() {
			case "int":
				i, _ := strconv.Atoi(name)
				v.FieldByName(f.Name).SetInt(int64(i))
			case "string":
				v.FieldByName(f.Name).SetString(name)
			case "float32":
			case "float64":
				ff, _ := strconv.ParseFloat(name, 10)
				v.FieldByName(f.Name).SetFloat(ff)
			}
		} else {
			return errors.New(fmt.Sprintf("parse request err, no param: %s", strings.ToLower(f.Name)))
		}
	}
	return nil
}

func buildTime(absoluteTime string) string {
	if len(absoluteTime) < 19 {
		return absoluteTime
	}
	abst := absoluteTime[:19]
	now := time.Now()
	target_time, err := time.ParseInLocation(TimeFormat, abst, LOCAL_LOCATION)
	get(err)
	sec := now.Unix() - target_time.Unix()

	hour := sec / 3600
	if hour == 0 {
		return fmt.Sprintf("%s|%d分钟前", abst, (sec / 60))
	}

	day := hour / 24
	if day == 0 {
		return fmt.Sprintf("%s|%d小时前", abst, hour)
	}
	if day == 7 {
		return fmt.Sprintf("%s|1周前", abst)
	}
	if day == 14 {
		return fmt.Sprintf("%s|半个月前", abst)
	}

	month := day / 30
	if month == 0 {
		return fmt.Sprintf("%s|%d天前", abst, day)
	}

	year := month / 12
	if year == 0 {
		return fmt.Sprintf("%s|%d个月前", abst, month)
	}

	return fmt.Sprintf("%s|%d年前", abst, year)
}

type price struct {
	Expire int64   `json:"expire"`
	Units  int64   `json:"units"`
	Money  float64 `json:"money"`
	Limit  int64   `json:"limit"`
}

func (p *price) chkParam() bool {
	if p.Expire <= 0 || p.Units < 0 || p.Money < 0 || p.Limit < 0 {
		return false
	}
	return true
}

func chkPrice(prices interface{}) *Error {
	b, err := json.Marshal(prices)
	if err != nil {
		return ErrParseJson(err)
	}

	pricePlans := []price{}

	json.Unmarshal(b, &pricePlans)
	if len(pricePlans) > DATAITEM_PRICE_MAX {
		return E(ErrorCodeItemPriceOutOfLimit)
	}
	if len(pricePlans) == 0 {
		return nil
	}

	for i, v := range pricePlans {
		if !v.chkParam() {
			return ErrInvalidParameter(fmt.Sprintf("price[%d]: %+v \n", i, v))
		}
	}

	return nil
}

func getPriceStat(prices interface{}) string {
	b, err := json.Marshal(prices)
	if err != nil {
		return ""
	}

	pricePlans := []price{}

	json.Unmarshal(b, &pricePlans)

	if len(pricePlans) == 0 {
		return ""
	}

	for _, v := range pricePlans {
		if v.Money == 0 {
			if v.Limit > 0 {
				return PRICE_STATE_FREE_LIMIT
			} else {
				return PRICE_STATE_FREE
			}
		} else {
			return PRICE_STATE_NOT_FREE
		}
	}

	return ""
}

func addPriceElemUid(price interface{}) {
	if arr, ok := price.([]interface{}); ok {
		for i, _ := range arr {
			if m, ok := arr[i].(map[string]interface{}); ok {
				if m["plan_id"] == nil {
					m["plan_id"] = bson.NewObjectId().Hex()
				}
			}
		}
	}
}

func (rep *repository) chkLabel() {
	if m, ok := rep.Label.(map[string]interface{}); ok {
		for _, v := range COL_LABEL_CHILDREN {
			if _, ok := m[v]; !ok {
				m[v] = make(map[string]interface{})
			}
		}
	} else {
		rep.Label = new(Label)
	}
}

func ifInLabel(i interface{}, column string) *Error {
	m, ok := i.(map[string]interface{})
	if !ok {
		return ErrNoParameter("label")
	}
	if mm := m["sys"]; mm != nil {
		if value := mm.(map[string]interface{})[column]; value != nil {
			if reflect.TypeOf(value).Kind() != reflect.String {
				return ErrNoParameter(fmt.Sprintf("label.sys.%s", column))
			} else if !contains(SUPPLY_STYLE_ALL, value.(string)) {
				return ErrInvalidParameter(fmt.Sprintf("label.sys.%s", column))
			}
		} else {
			return ErrNoParameter(fmt.Sprintf("label.sys.%s", column))
		}
	} else {
		return ErrNoParameter("label.sys")
	}

	for _, v := range COL_LABEL_CHILDREN {
		if _, ok := m[v]; !ok {
			m[v] = make(map[string]interface{})
		}
	}

	return nil
}

func contains(list interface{}, str string) bool {
	if list == nil {
		return false
	}

	switch list.(type) {

	case []interface{}:
		l, ok := list.([]interface{})
		if !ok {
			return false
		}
		for _, v := range l {
			if str == v.(string) {
				return true
			}
		}
	case []string:
		l, ok := list.([]string)
		if !ok {
			return false
		}
		for _, v := range l {
			if str == v {
				return true
			}
		}
	}
	return false
}

func httpGet(getUrl string, credential ...string) ([]byte, error) {
	var resp *http.Response
	var err error
	if len(credential) == 2 {
		req, err := http.NewRequest("GET", getUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("[http] err %s, %s\n", getUrl, err)
		}
		req.Header.Set(credential[0], credential[1])
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			Log.Error("http get err:%s", err.Error())
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
		}
	} else {
		resp, err = http.Get(getUrl)
		if err != nil {
			Log.Error("http get err:%s", err.Error())
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
		}
	}

	return ioutil.ReadAll(resp.Body)
}

func Env(name string, required bool, showLog ...bool) string {
	s := os.Getenv(name)
	if required && s == "" {
		panic("env variable required, " + name)
	}
	if len(showLog) == 0 || showLog[0] {
		Log.Infof("[env][%s] %s\n", name, s)
	}
	return s
}

func HttpPostJson(postUrl string, body []byte, credential ...string) ([]byte, error) {
	var resp *http.Response
	var err error
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", postUrl, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(credential[0], credential[1])
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", postUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http] status err %s, %d\n", postUrl, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s\n", postUrl, err)
	}
	return b, nil
}

func getToken(user, passwd string) string {
	passwdMd5 := getMd5(passwd)
	token := fmt.Sprintf("Basic %s", string(base64Encode([]byte(fmt.Sprintf("%s:%s", user, passwdMd5)))))
	URL := fmt.Sprintf("http://%s:%s", API_SERVER, API_PORT)
	b, err := httpGet(URL, AUTHORIZATION, token)
	if err != nil {
		Log.Errorf("get token err: %s", err.Error())
	}

	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		Log.Errorf("unmarshal token err: %s", err.Error())
	}

	return i.(map[string]interface{})["token"].(string)
}

func getMd5(content string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(content))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func Compare(src interface{}, dst interface{}) bool {
	v_src, v_dst := reflect.ValueOf(src), reflect.ValueOf(dst)

	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		v_src = v_src.Elem()
	}

	if reflect.TypeOf(dst).Kind() == reflect.Ptr {
		v_dst = v_dst.Elem()
	}

	if v_src.Kind() != reflect.Struct || v_dst.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v_src.NumField(); i++ {
		switch v_src.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v_src.Field(i).Int() != v_dst.Field(i).Int() {
				return false
			}
		case reflect.String:
			if v_src.Field(i).String() != v_dst.Field(i).String() {
				return false
			}
		case reflect.Float32, reflect.Float64:
			if v_src.Field(i).Float() != v_dst.Field(i).Float() {
				return false
			}
		case reflect.Bool:
			if v_src.Field(i).Bool() != v_dst.Field(i).Bool() {
				return false
			}
		default:
			Log.Info("compare unsupport type")
		}

	}

	return true
}

func getLabelValue(label interface{}, key string) interface{} {

	keys := strings.Split(key, ".")

	switch len(keys) {
	case 1:
		if l, ok := label.(map[string]interface{}); ok {
			return l[keys[0]]
		}
		if l, ok := label.(bson.M); ok {
			return l[keys[0]]
		}
	case 2:
		if l, ok := label.(map[string]interface{}); ok {
			return l[keys[0]].(map[string]interface{})[keys[1]].(string)
		}
		if l, ok := label.(bson.M); ok {
			return l[keys[0]].(bson.M)[keys[1]].(string)
		}
	}
	return nil
}

func getColumns(target interface{}) columns {
	t := reflect.TypeOf(target)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	n := t.NumField()
	cols := []column{}
	for i := 0; i < n; i++ {
		col := column{
			ColumName: t.Field(i).Tag.Get(DEFINE_TAG_NAME),
			ColumType: t.Field(i).Type.Name(),
		}
		cols = append(cols, col)
	}

	return columns(cols)
}

func cheParam(paramName, paramValue string) *Error {
	if strings.TrimSpace(paramValue) == "" {
		return ErrNoParameter(paramName)
	}

	var paramLengthLimit int
	switch paramName {
	case PARAM_ITEM_NAME:
		paramLengthLimit = LIMIT_ITEM_LENGTH
	case PARAM_REP_NAME:
		paramLengthLimit = LIMIT_REP_LENGTH
	case PARAM_TAG_NAME:
		paramLengthLimit = LIMIT_TAG_LENGTH
	}

	if len(paramValue) > paramLengthLimit {
		return ErrInvalidParameter(fmt.Sprintf("%s : %s", paramName, fmt.Sprintf("out of limit %d", paramLengthLimit)))
	}

	return nil
}

func ifCooperate(cooperate interface{}, loginName string) bool {
	return contains(cooperate, loginName)
}
