package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func putDataitemPermission(repname, itemname, username string) {
	selector := bson.M{
		COL_PERMIT_REPNAME:  repname,
		COL_PERMIT_ITEMNAME: itemname,
		COL_PERMIT_USER:     username,
	}
	update := Item_Permission{
		User_name:       username,
		Repository_name: repname,
		Dataitem_name:   itemname,
	}

	exec := Execute{
		Collection: C_DATAITEM_PERMISSION,
		Selector:   selector,
		Update:     update,
		Type:       Exec_Type_Upsert,
	}

	go asynExec(exec)

}

func putRepositoryPermission(repname, username string, opt_permission int) {
	if username == "" {
		return
	}
	selector := bson.M{
		COL_PERMIT_REPNAME: repname,
		COL_PERMIT_USER:    username,
	}
	update := Rep_Permission{
		User_name:       username,
		Repository_name: repname,
		Opt_permission:  opt_permission,
	}
	execs := []Execute{
		{
			Collection: C_REPOSITORY_PERMISSION,
			Selector:   selector,
			Update:     update,
			Type:       Exec_Type_Upsert,
		},
	}

	if update.Opt_permission == PERMISSION_WRITE {
		exec := Execute{
			Collection: C_REPOSITORY,
			Selector:   bson.M{COL_REPNAME: repname},
			Update:     bson.M{CMD_ADDTOSET: bson.M{COL_REP_COOPERATOR: username}},
			Type:       Exec_Type_Update,
		}
		execs = append(execs, exec)
	}
	go asynExec(execs...)
}

func getUserQuota(token, loginName string) Quota {
	b, err := httpGet(fmt.Sprintf("http://%s:%s/quota/%s/repository", API_SERVER, API_PORT, loginName), AUTHORIZATION, token)
	if err != nil {
		Log.Error(fmt.Sprintf("http://%s:%s/quota/%s/repository", API_SERVER, API_PORT, loginName), AUTHORIZATION, token)
		Log.Errorf("getUserQuota err :%s\n", err)
	}

	result := new(Result)
	err = json.Unmarshal(b, result)
	if err != nil {
		Log.Errorf("getUserQuota(%s) err :%s\n", loginName, err)
	}

	q := Quota{}
	if result.Data != nil {
		u := result.Data.(map[string]interface{})

		if pub, exist := u[Quota_Rep_Pub]; exist {
			q.Rep_Public, _ = strconv.Atoi(pub.(string))
		}
		if pri, exist := u[Quota_Rep_Pri]; exist {
			q.Rep_Private, _ = strconv.Atoi(pri.(string))
		}
		Log.Infof("user quota %#v\n", q)

		if q.Rep_Public == VIP_SERVICE_ADMIN_PUB {
			q.Rep_Public = 100000
		}
		if q.Rep_Private == VIP_SERVICE_ADMIN_PUB {
			q.Rep_Private = 100000
		}
	}

	return q
}

func getSubscribers(Type, repname, itemname, token string) []string {
	url := ""
	switch Type {
	case Subscripters_By_Rep:
		url = fmt.Sprintf("http://%s:%s/subscriptions/subscribers/%s?phase=1", API_SERVER, API_PORT, repname)
	case Subscripters_By_Item:
		url = fmt.Sprintf("http://%s:%s/subscriptions/subscribers/%s/%s?phase=1", API_SERVER, API_PORT, repname, itemname)
	}

	b, err := httpGet(url, AUTHORIZATION, token)
	if err != nil {
		Log.Errorf("get subscribers err :%s\n", err)
	}

	result := new(Result)
	err = json.Unmarshal(b, result)
	if err != nil {
		Log.Errorf("get subscribers err :%s\n", err)
	}

	return getResult(*result, "results")
}
