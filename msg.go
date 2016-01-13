package main

import (
	"encoding/json"
	"github.com/asiainfoLDP/datahub_repository/mq"
)

const (
	MQ_TOPIC_TO_SUB       = "to_subscriptions.json"
	MQ_TOPIC_TO_REP       = "to_repositories.json"
	MQ_KEY_ADD_PERMISSION = "add_permission"
	MQ_HANDLER_PERMISSION = "permission_handler"
	MQ_KEY                = "repositories"
)

type Msg struct {
	mq.MessageQueue
}

func (m *Msg) MqJson(topic string, content interface{}) {
	b, err := json.Marshal(content)
	get(err)
	msg.SendSyncMessage(topic, []byte(MQ_KEY), b)
}

type MyMesssageListener struct {
	name string
}

func newMyMesssageListener(name string) *MyMesssageListener {
	return &MyMesssageListener{name: name}
}

func (listener *MyMesssageListener) OnMessage(topic string, partition int32, offset int64, key, value []byte) bool {

	switch string(key) {
	case MQ_KEY_ADD_PERMISSION:
		m := make(Ms)
		if err := json.Unmarshal(value, &m); err != nil {
			Log.Errorf("%s received: (%d) message: %s", listener.name, offset, err.Error())
		}
		db.mqPermissionHandler(m)
	}
	return true
}

func (listener *MyMesssageListener) OnError(err error) bool {
	Log.Debugf("api response listener error: %s", err.Error())
	return false
}
