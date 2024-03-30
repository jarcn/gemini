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
	//var idMap = []int64{1,2,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,20,21,22,23,24,25,26,27,28,29,30,31,33,34,35,36,37,38,
	//	39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,55,56,57,59,60,61,64,67,68,69,70,71,72,73,74,76,78,79,80,81,82,83,84,
	//	86,87,88,89,90,93,96,97,98,99,100,101,102,104,105,106,107,108,109,110,111,112,113,115,117,118,120,122,124,125,126,
	//	127,128,129,131,134,135,136,138,139,140,141,142,144,145,146,148,149,150,152,153,154,155,156,157,158,159,160,161,
	//	162,163,164,165,166,167,168,169,170,171,172,173,174,176,178,179,181,182,183,185,186,187,188,189,190,191,192,193,
	//	194,195,196,197,199,200,203,204,205,207,208,209,210,211,212,213,214,215,216,217,218,220,222,225,226,227,228,229,
	//	230,231,233,234,235,236,237,238,239,240,241,242}
	var idMap = []int64{1, 2, 4, 5}
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
