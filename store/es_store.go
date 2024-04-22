package store

import (
	"bytes"
	"context"
	"encoding/json"
	"gemini/profile"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strconv"
)

func Insert2ES(resumeArr []profile.Resume) {
	// Elasticsearch 连接配置
	cfg := elasticsearch.Config{
		//Addresses: []string{"http://10.128.0.165:9200", "http://10.128.0.72:9200"}, //生产环境
		Addresses: []string{"http://10.129.0.251:9200", "http://10.129.0.217:9200", "http://10.129.0.146:9200"}, //预发环境
		Password:  "Qiyi123!@#",
		Username:  "elastic",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 执行批量插入请求
	// 批量插入数据
	for _, doc := range resumeArr {
		// 将文档转换为 JSON 格式
		docJSON, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling document: %s", err)
			continue
		}

		// 准备批量插入请求
		req := esapi.IndexRequest{
			//Index:      "gemini_hra", // 替换为你的索引名称
			Index:      "hra_cv", // 替换为你的索引名称
			DocumentID: strconv.Itoa(int(doc.ID)),
			Body:       bytes.NewReader(docJSON),
			Refresh:    "true",
		}

		// 执行批量插入请求
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("Error indexing document: %s", err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document: %s", res.Status(), res.String())
		} else {
			log.Printf("[%s] Indexed document with ID: %d", res.Status(), doc.ID)
		}
	}
}

func SelectFrom2ES(resumeArr []profile.Resume) {
	// Elasticsearch 连接配置
	cfg := elasticsearch.Config{
		Addresses: []string{"http://10.129.0.251:9200", "http://10.129.0.217:9200", "http://10.129.0.146:9200"}, //预发环境
		Password:  "Qiyi123!@#",
		Username:  "elastic",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 执行批量插入请求
	// 批量插入数据
	for _, doc := range resumeArr {
		// 将文档转换为 JSON 格式
		docJSON, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling document: %s", err)
			continue
		}

		// 准备批量插入请求
		req := esapi.IndexRequest{
			//Index:      "gemini_hra", // 替换为你的索引名称
			Index:      "hra_cv", // 替换为你的索引名称
			DocumentID: strconv.Itoa(int(doc.ID)),
			Body:       bytes.NewReader(docJSON),
			Refresh:    "true",
		}

		// 执行批量插入请求
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("Error indexing document: %s", err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document: %s", res.Status(), res.String())
		} else {
			log.Printf("[%s] Indexed document with ID: %d", res.Status(), doc.ID)
		}
	}
}
