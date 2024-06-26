package main

import (
	"encoding/json"
	"gemini/db"
	"gemini/store"
	"gemini/tasks"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
	//db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //生产环境
	//cache.InitKeyCache()
}

// var keys = []string{"AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90", "AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90", "AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90"}
var keys = []string{"AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90"}

func main() {
	// 创建channel用于存储数据
	dataChan := make(chan store.GltZhData, 1000)
	// 查询所有数据并存储到channel中
	go func() {
		var d = store.GltZhData{}
		for start := 0; ; start += 3000 {
			end := start + 3000
			allData, _ := d.SelectByPage(start, end, db.Client())
			if len(allData) == 0 {
				close(dataChan)
				break
			}
			for _, datum := range allData {
				dataChan <- datum
			}
		}
	}()

	// 启动4个goroutines并行消费数据
	for i := 0; i < 4; i++ {
		go func(key string) {
			for datum := range dataChan {
				data := map[string]string{
					"url":     datum.URL,
					"profile": datum.Profile,
				}
				marshal, _ := json.Marshal(data)
				merge := tasks.SyncDoMerge(marshal, key)
				if merge != nil {
					sd := map[string]int64{
						"id": merge.ID,
					}
					step2Json, _ := json.Marshal(sd)
					tasks.DoDeduce(step2Json, key, false)
				}
			}
		}(keys[i%len(keys)]) // 循环使用keys中的key
	}

	// 等待所有goroutines完成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gemini merge server...")
	db.Close()
	log.Println("Gemini merge server stopped")
}
