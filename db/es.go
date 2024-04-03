package db

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

var esClient *elasticsearch.Client

func MustInitES(adds []string, userName, passWord string) {
	cfg := elasticsearch.Config{
		Addresses: adds,
		Password:  userName,
		Username:  passWord,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	esClient = es
}

func EsInsert(indexName string, profiles []profile.Resume) {
	// 批量插入数据
	for _, doc := range profiles {
		// 将文档转换为 JSON 格式
		docJSON, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling document: %s", err)
			continue
		}
		// 准备批量插入请求
		req := esapi.IndexRequest{
			Index:      indexName,
			DocumentID: strconv.Itoa(int(doc.ID)),
			Body:       bytes.NewReader(docJSON),
			Refresh:    "true",
		}
		// 执行批量插入请求
		res, err := req.Do(context.Background(), esClient)
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
