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
	//db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //生产环境
	//cache.InitKeyCache()
}

var keys = []string{"AIzaSyD50ffX7kVQs7AYAR-MBGLCs5O_LxCKOfQ", "AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90", "AIzaSyD50ffX7kVQs7AYAR-MBGLCs5O_LxCKOfQ", "AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90", "AIzaSyD50ffX7kVQs7AYAR-MBGLCs5O_LxCKOfQ", "AIzaSyB5Yx-nRni3ILCiD8CAc8zKTOcEFInDv90"}

func main() {
	for i := 0; i < len(keys); i++ {
		go run(keys[i], i*2000, 2000)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gemini merge server...")
	db.Close()
	log.Println("Gemini merge server stopped")
}

func run(key string, start, end int) {
	var d = store.GltZhData{}
	allData, _ := d.SelectByPage(start, end, db.Client())
	log.Printf("start:%d end:%d query %d rows,start do task...\r\n", start, end, len(allData))
	for _, datum := range allData {
		data := make(map[string]string)
		data["url"] = datum.URL
		data["profile"] = datum.Profile
		marshal, _ := json.Marshal(data)
		merge := tasks.SyncDoMerge(marshal, key)
		if merge != nil {
			sd := make(map[string]int64)
			sd["id"] = merge.ID
			step2Json, _ := json.Marshal(sd)
			tasks.DoDeduce(step2Json, key, false)
		}
	}
}
