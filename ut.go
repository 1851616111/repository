package main

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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
}
