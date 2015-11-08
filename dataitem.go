package main

import (
	"fmt"
)

func (db *DB) getDataitems(columnName, columnValue string) ([]DataItem, error) {
//	item := new(DataItem)
//	rows, err := db.Where(fmt.Sprintf(" %s = ? ", columnName), columnValue).Rows(item)
//	if err != nil {
//		log.Printf("getDataitem err", err)
//		return nil, err
//	}
//	defer rows.Close()
//	res := []DataItem{}
//	for rows.Next() {
//		i := new(DataItem)
//		if err := rows.Scan(i); err != nil {
//			return nil, err
//		}
//		res = append(res, *i)
//	}
//	return res, nil
	l := []DataItem{}
	var err error
	if columnName == "" || columnValue == ""  {
		err = db.Find(&l)
	} else {
		err = db.Where(fmt.Sprintf(" %s = ? ", columnName), columnValue).Find(&l)
	}

	return l, err


}
func (db *DB) getDataitemByIds(dataitemIds ...interface{}) (map[int64]DataItem, error) {
	m := map[int64]DataItem{}
	err := db.In("DATAITEM_ID", dataitemIds...).Find(&m)
	return m, err
}

func (db *DB) getDataitem(d *DataItem) (bool, error) {
	return db.Get(d)
}

func (db *DB) setDataitem(d *DataItem) error {
	_, err := db.InsertOne(d)
	return err
}

func (db *DB) getDataitem_Chosen(chosenName ...string) ([]Dataitem_Chosen, error) {
	l := []Dataitem_Chosen{}
	var err error
	if len(chosenName) == 0 {
		err = db.Find(&l)
	} else {
		err = db.Where(" CHOSEN_NAME = ? ", chosenName[0]).Find(&l)
	}
	return l, err
}
func (db *DB) getDataitem_ChosenNames() ([]Dataitem_Chosen, error) {
	l := []Dataitem_Chosen{}
	err := db.Distinct("CHOSEN_NAME").Find(&l)
	return l, err
}

func (db *DB) setDataitem_Chosen(d *Dataitem_Chosen) error {
	_, err := db.InsertOne(d)
	return err
}

func (db *DB) getDataitemsByIds(ids []interface{}) []Data {
	l := []Data{}
	if len(ids) == 0 {
		return l
	}
	usages_m, err := db.getDataitemUsageByIds(ids...)
	get(err)
	items_m, err := db.getDataitemByIds(ids...)
	get(err)

	m := make(M)
	for k, item := range items_m {
		usage := usages_m[k]
		usage.Refresh_date = buildTime(usage.Refresh_date)
		d := Data{Item: &item, Usage: &usage}
		m[k] = d
	}

	for _, v := range m {
		l = append(l, v.(Data))
	}
	return l
}
