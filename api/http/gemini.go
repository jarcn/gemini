package http

import (
	"gemini/api/bean"
	"gemini/queue"
	"gemini/workers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func FormatCV(c *gin.Context, pool *queue.GeminiQueue) {
	var data bean.SourceData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "failed",
		})
		log.Fatal(err)
		return
	}
	// todo 接收到数据进行 gemini 调用
	workers.CallGemini(data.ResumeData, pool)
	c.JSON(http.StatusOK, gin.H{
		"code": "ok",
	})
}
