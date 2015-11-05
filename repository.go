package main

import (
	"fmt"
	"log"
)

func (db *DB) getRepository(columnName, columnValue string) []Repository {
	repo := new(Repository)
	rows, err := db.Where(fmt.Sprintf(" %s = ? ", columnName), columnValue).Rows(repo)
	if err != nil {
		log.Printf("getRepository err", err)
	}
	defer rows.Close()
	res := []Repository{}
	for rows.Next() {
		err = rows.Scan(repo)
		res = append(res, *repo)
	}
	return res
}
