package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
		log.Println(err)
	}
}

type Result struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

//w.Header().Set("Content-Type", "application/json; charset=utf-8")
//w.WriteHeader(statusCode)
//
//if e == nil {
//e = ErrorNone
//}
//
////result := Result {code: e.code, msg: e.message, data: data}
//result := Result {Code: e.code, Msg: e.message}
//
//log.Printf (fmt.Sprintf ("code: %s", e.code))
//
//jsondata, err := json.Marshal(&result)
//if err != nil {
//w.Write([]byte(getJsonBuildingErrorJson ()))
//} else {
//w.Write(jsondata)
//}

func (p *Rsp) Json(code int, e *Error, data ...interface{}) (int, string) {
	p.w.Header().Set("Access-Control-Allow-Origin", "*")
	p.w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	p.w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept,X-Requested-With")
	p.w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	result := Result{Code: e.code, Msg: e.message, Data: data}
	b, err := json.Marshal(result)
	chk(err)
	return code, string(b)
}

func Env(name string, required bool) string {
	s := os.Getenv(name)
	if required && s == "" {
		panic("env variable required, " + name)
	}
	return s
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

func buildTime(absoluteTime string) string {
	now := time.Now()

	sevenDayAgo := now.AddDate(0, 0, -7)
	target_time, err := time.ParseInLocation(TimeFormat, absoluteTime, LOCAL_LOCATION)
	get(err)
	if target_time.After(sevenDayAgo) {
		sec := now.Unix() - target_time.Unix()

		oneDayAgo := now.AddDate(0, 0, -1)

		if target_time.After(oneDayAgo) {
			return fmt.Sprintf("%s,%d小时以前", absoluteTime, sec/3600)
		} else {
			return fmt.Sprintf("%s,%d天以前", absoluteTime, sec/(3600*24))
		}
	}
	return absoluteTime
}

func (p *repository) ParseRequeset(r *http.Request) error {
	t := reflect.TypeOf(*p)
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name := strings.TrimSpace(r.PostFormValue(strings.ToLower(f.Name))); name != "" {
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
			case "bool":
				b, err := strconv.ParseBool(name)
				chk(err)
				v.FieldByName(f.Name).SetBool(b)
			}
		}
	}
	return nil
}

func (p *repository) BuildRequest() {
	if p.Repaccesstype == "" {
		p.Repaccesstype = "public"
	}

	p.Optime = time.Now()
}

func (p *dataItem) ParseRequeset(r *http.Request) {
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
			case "bool":
				b, err := strconv.ParseBool(name)
				chk(err)
				v.FieldByName(f.Name).SetBool(b)
			}
		}
	}
}

func (p *dataItem) BuildRequeset(repName, itemName, createName string) {
	p.Repository_name = repName
	p.Dataitem_name = itemName
	p.Create_name = createName
	p.Optime = time.Now()
	if p.Itemaccesstype == "" {
		p.Itemaccesstype = "public"
	}
	p.Ct = time.Now()
}
