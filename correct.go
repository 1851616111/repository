package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func correctQuota(db *DB) {
	copy := db.copy()
	defer copy.Close()
	aggregate := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"create_user": "$create_user",
					"accesstype":  "$repaccesstype",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	var correct func(m Ms, token string) = func(m Ms, token string) {
		if len(m) > 0 {
			if mm, ok := m["_id"].(Ms); ok {
				username, accesstype, count := mm["create_user"], mm["accesstype"], m["count"]
				context := fmt.Sprintf(`{"%s":%d}`, accesstype, count)
				url := fmt.Sprintf("http://%s:%s/quota/%s/repository/use", API_SERVER, API_PORT, username)
				res, err := HttpPostJson(url, []byte(context), AUTHORIZATION, token)
				if err != nil {
					Log.Errorf("-----ERR %s", err.Error())
				}
				Log.Error("------------> res%#v", res)
			}
		}
	}

	token := getToken(Username, Password)
	Log.Infof("init token ok. token: %s", token)
	result := Ms{}
	pipe := copy.DB(DB_NAME).C(C_REPOSITORY).Pipe(aggregate)
	iter := pipe.Iter()
	for iter.Next(&result) {
		correct(result, token)
	}
	if err := iter.Close(); err != nil {
		Log.Errorf("-----ERR %s", err.Error())
	}
}
