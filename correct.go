package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func correctQuota(db *DB) {

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
				_, err := HttpPostJson(url, []byte(context), AUTHORIZATION, token)
				if err != nil {
					Log.Errorf("update user quota err %s", err.Error())
				}
			}
		}
	}

	fn := func(db *DB, token string) {
		copy := db.copy()
		defer copy.Close()

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

	token := getToken(Username, Password)
	Log.Infof("init token ok. token: %s", token)

	for {
		fn(db, token)
		time.Sleep(time.Minute * 60)
	}
}
