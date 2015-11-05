package main

import (
	"github.com/lunny/log"
)

func userRepositoryPermit(user_id string, db *DB) []string {
	permit := new(Repository_Permit)
	rows, err := db.Cols("repository_name").Where(" user_id = ? ", user_id).Rows(permit)
	if err != nil {
		log.Printf("userRepositoryPermit err", err)
	}
	defer rows.Close()
	res := []string{}
	for rows.Next() {
		err = rows.Scan(permit)
		res = append(res, permit.Repository_name)
	}
	return res
}
