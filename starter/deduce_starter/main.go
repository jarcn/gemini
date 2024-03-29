package main

import (
	"encoding/json"
	"fmt"
	"gemini/cache"
	"gemini/db"
	"gemini/store"
	"gemini/tasks"
	"time"
)

func init() {
	//db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //预发环境
	cache.InitKeyCache()
}

func main() {
	var idMap = []int64{1, 2, 3, 4, 5}
	for _, id := range idMap {
		m1 := make(map[string]int64)
		m1["id"] = id
		marshal, _ := json.Marshal(m1)
		deduce := tasks.DoDeduce(marshal, "AIzaSyATzYCPsgwDbzXdrcF5V5AKom3hr8MZwZ4")
		fmt.Println(deduce)
		fmt.Println("sleep 30s")
		time.Sleep(time.Second * 30)
	}
}
func getKey() string {
	key := cache.GetKey()
	result := store.GeminiResult{}
	count, err := result.CountByKey(db.Client(), key)
	if err != nil {
		return cache.GetKey()
	}
	if count == 0 {
		return key
	}
	currentTime := time.Now().Unix()
	if currentTime-count > 60 {
		return key
	} else {
		time.Sleep(time.Second * 60)
	}
	return cache.GetKey()
}
