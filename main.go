package main

import (
	"gemini/cache"
	"gemini/db"
	"gemini/tasks"
	"log"
)

func init() {
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data")
	cache.InitKeyCache()
}

func main() {
	// 数据处理程序已经准备好
	// TODO 组织数据进行处理
	result := tasks.CallGemini("", "", cache.GetKey())
	log.Println(result)
}
