package models

import (
	"database/sql"
	mysql "gemini/store/client"
	"log"
)

type Event struct {
	Id         int64         `db:"id" json:"Id"`
	EventName  string        `db:"event_name" json:"EventName"` //事件名称
	ThirdType  int           `db:"third_type" json:"ThirdType"` //0:所有三方;1:Moengage;2:Appsflyer
	CreateTime sql.NullInt64 `db:"create_time,omitempty" json:"CreateTime,omitempty"`
	UpdateTime sql.NullInt64 `db:"update_time,omitempty" json:"UpdateTime,omitempty"`
}

func selectAll() *[]Event {
	var events []Event
	all := `select * from qiyee_job_system.tbl_sys_third_event`
	err := mysql.DbClient().Select(&events, all)
	if err != nil {
		log.Panic("select tbl_sys_third_event data error", err)
	}
	return &events
}

var thirdEventIdCache = make(map[string]int, 4096)

func GetEventInfo(eventName string) int {
	if value, ok := thirdEventIdCache[eventName]; ok {
		return value
	} else {
		thirdEventIdCache[eventName] = -1 //todo cache 需优化
		refreshCache()
		if v, k := thirdEventIdCache[eventName]; k {
			return v
		} else {
			return -1
		}
	}
}

func refreshCache() {
	allEvent := selectAll()
	for _, event := range *allEvent {
		thirdEventIdCache[event.EventName] = event.ThirdType
	}
	log.Println("third event id refresh success")
}
