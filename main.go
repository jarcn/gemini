package main

import (
	"flag"
	"gemini/api/route"
	"gemini/configs"
	"gemini/queue"
	"gemini/store/client"
	"github.com/gin-gonic/gin"
	"strings"
)

var ServerConfig configs.Config
var GeminiPool *queue.GeminiQueue

func init() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件地址")
	flag.Parse()
	if configPath == "" {
		panic("You must have one config.")
	}
	ServerConfig = configs.GetConfig(configPath)
}

func init() {
	if ServerConfig.SourceKey != "" {
		GeminiPool = queue.NewSafeQueue()
		keys := strings.Split(ServerConfig.SourceKey, ",")
		queue.InitQueue(keys, GeminiPool)
	}
}

func main() {
	r := gin.Default()
	route.Register(r, GeminiPool)
	defer func() {
		client.Close() //关闭mysql session
	}()
	r.Run(":8080")
}
