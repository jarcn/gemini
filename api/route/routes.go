package route

import (
	"gemini/api/http"
	"gemini/queue"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, pool *queue.GeminiQueue) {
	// 初始页面
	r.GET("/", http.Welcome)

	r.POST("/convert", func(context *gin.Context) {
		http.FormatCV(context, pool)
	})
}
