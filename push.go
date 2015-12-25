package main

import (
	"fmt"
	"time"
)

const (
	Item_Channel_Max = 10000
	Service_Name     = "repository"
	PayLoad_Table    = "dataitem"
)

var (
	m_Cache_Rep = make(map[string]meta_rep)
	cols        columns
)

func init() {
	cols = getColumns(meta_item{})
}

func pushMetaDataLoop(db *DB) {
	timer1 := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-timer1.C:
			if time.Now().Hour() == 23 {
				copy := db.copy()
				pushMetaData(copy, &msg)
			}
		}
	}
}

func pushMetaData(src *DB, dst *Msg) {
	defer src.Close()
	reps, _ := src.getRepositories(nil)
	go func(reps *[]repository) {
		for _, v := range *reps {
			meta := meta_rep{
				User:       v.Create_user,
				Rep:        v.Repository_name,
				AccessType: v.Repaccesstype,
			}
			m_Cache_Rep[meta.Rep] = meta
		}
	}(&reps)

	items, _ := db.getDataitems(0, SELECT_ALL, nil)
	go func(items *[]dataItem, dst *Msg) {
		data := []meta_item{}

		for i, v := range *items {
			rep := m_Cache_Rep[v.Repository_name]

			meta := meta_item{
				Id:             i,
				RepUser:        rep.User,
				Rep:            v.Repository_name,
				RepAccessType:  rep.AccessType,
				ItemUser:       v.Create_user,
				Item:           v.Dataitem_name,
				Tags:           v.Tags,
				ItemAccessType: v.Itemaccesstype,
			}

			suppleStyle, selectLabel := getLabelValue(v.Label, fmt.Sprintf("sys.%s", COL_ITEM_SYPPLY_STYLE)), getLabelValue(v.Label, fmt.Sprintf("sys.%s", COL_SELECT_LABEL))
			if suppleStyle != nil {
				if s, ok := suppleStyle.(string); ok {
					meta.SuppleStyle = s
				}
			}
			if selectLabel != nil {
				if s, ok := selectLabel.(string); ok {
					meta.SelectLabel = s
				}
			}

			data = append(data, meta)
		}

		pld := payLoad{
			Service: Service_Name,
			Date:    time.Now().Format(TimeFormatDay),
			Columns: cols,
			Table:   PayLoad_Table,
			Data:    data,
		}

		dst.MqJson(pld)

	}(&items, dst)

}

type meta_rep struct {
	User       string
	Rep        string
	AccessType string
}

type meta_item struct {
	Id             int    `json:"id" column:"repuser"`
	RepUser        string ` json:"repuser" column:"repuser"`
	Rep            string ` json:"repname" column:"repname"`
	RepAccessType  string ` json:"repaccesstype" column:"repaccesstype"`
	ItemUser       string ` json:"itemuser" column:"itemuser"`
	Item           string ` json:"itemname" column:"itemname"`
	Tags           int    ` json:"tags" column:"tags"`
	ItemAccessType string ` json:"itemaccesstype" column:"itemaccesstype"`
	SuppleStyle    string ` json:"supplestyle" column:"supplestyle"`
	SelectLabel    string ` json:"selectlabel" column:"selectlabel"`
}

type payLoad struct {
	Service string      `json:"service"`
	Date    string      `json:"date"`
	Columns columns     `json:"columns"`
	Table   string      `json:"table"`
	Data    interface{} `json:"data"`
}

type columns []column

type column struct {
	ColumName string `json:"column_name"`
	ColumType string `json:"column_type"`
}
