package main

import (
	"gopkg.in/mgo.v2/bson"
)

func setDataitemPermission(repname, itemname, username string) {
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
