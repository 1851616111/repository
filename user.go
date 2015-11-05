package main

import "fmt"

func (db *DB) userExist(user_id, passwd string) (bool, string) {
	user := new(User)
	has, err := db.Cols("user_id").Where(" EMAIL = ? and LOGIN_PASSWD = ?", user_id, passwd).Get(user)
	get(err)
	return has, fmt.Sprintf("%d", user.User_id)
}
