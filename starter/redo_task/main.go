package main

import (
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
	dataChan := make(chan store.GeminiResult, 1000)
	// 查询所有数据并存储到channel中
	go func() {
		var d = store.GeminiResult{}
		allData, _ := d.FindAllError(db.Client())
		for _, datum := range allData {
			dataChan <- *datum
		}
	}()
	// 启动4个goroutines并行消费数据
	for i := 0; i < 4; i++ {
		go func(key string) {
			for datum := range dataChan {
				merge := tasks.ReDoMergeBeat(datum, key)
				if merge != nil {
					tasks.ReDoDeduceBeat(*merge, key)
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
