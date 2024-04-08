package main

import (
	"gemini/cache"
	"gemini/db"
	"gemini/starter/kfk_consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 数据库
func initMysql() {
	db.MustInitMySQL(db.MysqlAddr)
	cache.InitKeyCache()
}

// ES
func initEs() {
	db.MustInitES(db.EsBrokers, "Qiyi123!@#", "elastic")
}

func main() {
	initMysql()
	initEs()
	go kfk_consumer.Step1ConsumerStart(db.KafkaBrokers, db.Step1Topic, db.Step1Group)
	go kfk_consumer.Step2ConsumerStart(db.KafkaBrokers, db.Step2Topic, db.Step2Group)
	log.Println("gemini merge server start success")
	// 保持服务持续运行
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gemini merge server...")
	db.Close()
	log.Println("Gemini merge server stopped")
}
