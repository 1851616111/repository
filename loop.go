package main

import (
	"fmt"
	"log"
	"time"
)

func DimLoop(db *DB) {
	for {
		updateDim(db)
		time.Sleep(time.Hour * 24)
	}
}

func updateDim(db *DB) {
	l := []Dim_Table{}
	db.Cols("FIELD_NAME", "ID", "NAME").Find(&l)

	for _, v := range l {
		_, exists := dim.GetM(v.Field_name)
		if !exists {
			m := make(M)
			dim.SetM(v.Field_name, m)
		}
		dim.Set(v.Field_name, fmt.Sprintf("%d", v.Id), v.Name)
	}
	log.Println("init dim table over")
}
