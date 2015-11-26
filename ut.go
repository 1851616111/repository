package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	TimeFormat = "2006-01-02 15:04:05"
)

var LOCAL_LOCATION *time.Location

func init() {
	loc, err := time.LoadLocation("Local")
	chk(err)
	LOCAL_LOCATION = loc
}

type MM map[interface{}]M
type M map[interface{}]interface{}
type Ms map[string]interface{}

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
	abst := absoluteTime[:len(absoluteTime)-10]
	now := time.Now()
	sevenDayAgo := now.AddDate(0, 0, -7)
	target_time, err := time.ParseInLocation(TimeFormat, abst, LOCAL_LOCATION)
	get(err)
	if target_time.After(sevenDayAgo) {
		sec := now.Unix() - target_time.Unix()
		oneDayAgo := now.AddDate(0, 0, -1)
		if target_time.After(oneDayAgo) {
			hour := sec / 3600
			if hour == 0 {
				return fmt.Sprintf("%s|%d分钟以前", abst, (sec%3600)/60)
			}
			return fmt.Sprintf("%s|%d小时以前", abst, hour)
		} else {
			return fmt.Sprintf("%s|%d天以前", abst, sec/(3600*24))
		}
	}
	return abst
}

//func Parse(p *dataItem, cb ...func()) {
//	t := reflect.TypeOf(*p)
//	v := reflect.ValueOf(p).Elem()
//	for i := 0; i < t.NumField(); i++ {
//		f := t.Field(i)
//		if name := r.PostFormValue(strings.ToLower(f.Name)); name != "" {
//			switch f.Type.Name() {
//			case "int":
//				i, _ := strconv.Atoi(name)
//				v.FieldByName(f.Name).SetInt(int64(i))
//			case "string":
//				v.FieldByName(f.Name).SetString(name)
//			case "float32":
//			case "float64":
//				ff, _ := strconv.ParseFloat(name, 10)
//				v.FieldByName(f.Name).SetFloat(ff)
//			}
//		}
//	}
//}

//func paramSetter(p *dataItem, cb ...func()) {
//	t := reflect.TypeOf(*p)
//	v := reflect.ValueOf(p).Elem()
//	for i := 0; i < t.NumField(); i++ {
//		f := t.Field(i)
//		switch f.Type.Name() {
//		case "int":
//			i, _ := strconv.Atoi(name)
//			v.FieldByName(f.Name).SetInt(int64(i))
//		case "string":
//			v.FieldByName(f.Name).SetString(name)
//		case "float32":
//		case "float64":
//			ff, _ := strconv.ParseFloat(name, 10)
//			v.FieldByName(f.Name).SetFloat(ff)
//		case "interface{}":
//
//		}
//		v.Field(i).
//		if name := r.PostFormValue(strings.ToLower(f.Name)); name != "" {
//			switch f.Type.Name() {
//			case "int":
//				i, _ := strconv.Atoi(name)
//				v.FieldByName(f.Name).SetInt(int64(i))
//			case "string":
//				v.FieldByName(f.Name).SetString(name)
//			case "float32":
//			case "float64":
//				ff, _ := strconv.ParseFloat(name, 10)
//				v.FieldByName(f.Name).SetFloat(ff)
//			}
//	}
//}

func ifInLabel(i interface{}, column string) *Error {
	if m := i.(map[string]interface{})["sys"]; m != nil {
		if value := m.(map[string]interface{})[column]; value != nil {
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
	return nil
}

func contains(l []string, str string) bool {
	for _, v := range l {
		if str == v {
			return true
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
	} else {
		resp, err = http.Get(getUrl)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
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
