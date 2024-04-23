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
	db.MustInitMySQL("root:qiyi123!@#@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
}

var keys = []string{"AIzaSyA7TVlO6k72aC1SWrZBPHYuqV04NYPGw8I"}

func main() {
	// 创建channel用于存储数据
	dataChan := make(chan store.HraCvData, 1000)
	// 查询所有数据并存储到channel中
	go func() {
		var d = store.HraCvData{}
		allData, _ := d.SelectAll(db.Client())
		log.Printf("query total %d rows data \r\n", len(allData))
		for _, datum := range allData {
			dataChan <- datum
		}
	}()

	// 启动4个goroutines并行消费数据
	for i := 0; i < 4; i++ {
		go func(key string) {
			for datum := range dataChan {
				data := map[string]string{
					"url":     datum.ResumeLink,
					"profile": datum.Content,
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
		}(keys[0]) // 循环使用keys中的key
	}

	// 等待所有goroutines完成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gemini merge server...")
	db.Close()
	log.Println("Gemini merge server stopped")
}
