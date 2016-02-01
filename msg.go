package main

import (
	"encoding/json"
	"github.com/asiainfoLDP/datahub_repository/mq"
)

const (
	MQ_TOPIC_TO_SUB             = "to_subscriptions.json"
	MQ_TOPIC_TO_REP             = "to_repositories.json"
	MQ_TOPIC_TO_STATIS          = "statistics"
	MQ_KEY_ADD_PERMISSION       = "add_permission"
	MQ_KEY_ADD_STATIS_RANK_REP  = "add_statis_rank_rep"
	MQ_KEY_ADD_STATIS_RANK_ITEM = "add_statis_rank_item"
	MQ_HANDLER_PERMISSION       = "permission_handler"
	MQ_KEY                      = "repositories"
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

	if len(value) > 0 {
		Log.Infof("[DEBUG]------------- get msg %s ------------------>\n", string(value))

		switch string(key) {
		case MQ_KEY_ADD_PERMISSION:
			m := make(map[string]interface{})
			if err := json.Unmarshal(value, &m); err != nil {
				Log.Errorf("%s received: (%d) message: %s", listener.name, offset, err.Error())
			}
			mqPermissionHandler(m)

		case MQ_KEY_ADD_STATIS_RANK_REP:

			result := []statisRepRank{}
			if err := json.Unmarshal(value, &result); err != nil {
				Log.Errorf("%s received: (%d) message: %s", listener.name, offset, err.Error())
			}
			mqRankHandler(result)

		case MQ_KEY_ADD_STATIS_RANK_ITEM:

			result := []statisItemRank{}
			if err := json.Unmarshal(value, &result); err != nil {
				Log.Errorf("%s received: (%d) message: %s", listener.name, offset, err.Error())
			}
			mqRankHandler(result)
		}
	}
	return true
}

func (listener *MyMesssageListener) OnError(err error) bool {
	Log.Debugf("api response listener error: %s", err.Error())
	return false
}

type statisRepRank struct {
	Repository_name string
	Rank            float64
}

type statisItemRank struct {
	Repository_name string
	Dataitem_name   string
	Rank            float64
}
