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
	go kfk_consumer.Step1ConsumerStart(db.KafkaBrokers, "step1-gemini-merge", "step1-gemini-group")
	go kfk_consumer.Step2ConsumerStart(db.KafkaBrokers, "step2-gemini-deduce", "step2-gemini-group")
	log.Println("gemini merge server start success")
	// 保持服务持续运行
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gemini merge server...")
	db.Close()
	log.Println("Gemini merge server stopped")
}
