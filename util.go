package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
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

func (p *Rsp) Json(code int, data interface{}) (int, string) {
	b, err := json.Marshal(data)
	chk(err)
	p.w.Header().Set("Access-Control-Allow-Origin", "*")

	p.w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	p.w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept,X-Requested-With")
	p.w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	p.w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return code, string(b)
}
func str2ui64(p string) interface{} {
	if p == "" {
		return uint64(0)
	}
	i, err := strconv.ParseInt(p, 10, 64)
	chk(err)
	return uint64(i)
}
func str2ui8(p string) interface{} {
	if p == "" {
		return uint8(0)
	}
	i, err := strconv.ParseInt(p, 10, 8)
	chk(err)
	return uint8(i)
}
func str2f64(p string) interface{} {
	if p == "" {
		return 0.0
	}
	i, err := strconv.ParseFloat(p, 64)
	chk(err)
	return i
}

func random4Num() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(1000)
}
func Env(name string, required bool) string {
	s := os.Getenv(name)
	if required && s == "" {
		panic("env variable required, " + name)
	}
	return s
}

func (p *DataItem) ParseRequeset(r *http.Request) {
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
		}
	}
}

func (p *Dataitem_Chosen) ParseRequeset(r *http.Request) error {
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
func (p *DataItem) BuildRequeset(repName, itemName, user_id string) {
	p.Repository_name = repName
	p.Dataitem_name = itemName
	i, _ := strconv.Atoi(user_id)
	p.User_id = i
	p.Optime = time.Now().Format("2006-01-02 15:04:05")
}

func md5Passed(passwd string) string {
	h := md5.New()
	h.Write([]byte(passwd))
	return hex.EncodeToString(h.Sum(nil))
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
