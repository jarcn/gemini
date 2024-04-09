package main

import (
	"encoding/json"
	"gemini/cache"
	"gemini/db"
	"gemini/store"
	"gemini/tasks"
)

func init() {
	db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
	cache.InitKeyCache()
}

func main() {
	var d = store.GltZhData{}
	allData, _ := d.SelectAllData(db.Client())
	for _, datum := range allData {
		data := make(map[string]string)
		data["url"] = datum.URL
		data["profile"] = datum.Profile
		marshal, _ := json.Marshal(data)
		merge := tasks.SyncDoMerge(marshal, cache.GetKey())
		if merge != nil {
			sd := make(map[string]int64)
			sd["id"] = merge.ID
			step2Json, _ := json.Marshal(sd)
			tasks.DoDeduce(step2Json, cache.GetKey(), false)
		}
	}
}
