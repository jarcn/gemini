package main

import (
	"bytes"
	"context"
	"encoding/json"
	"gemini/db"
	"gemini/profile"
	"gemini/store"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strconv"
)

func init() {
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //预发环境
	//db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
}

func main() {
	result := store.GeminiResult{}
	all, _ := result.FindAll(db.Client())
	var resumeArr []profile.Resume
	for _, d := range all {
		step1 := d.GeminiStep1
		step2 := d.GeminiStep2
		resume := profile.MergeStep1AndStep2([]byte(step1), []byte(step2))
		resume.BasicInformation.ProfileUrl = d.CVURL
		resume.ID = d.ID
		resumeArr = append(resumeArr, *resume)
	}
	insert2ES(resumeArr)
}

func insert2ES(resumeArr []profile.Resume) {
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
			Index:      "hra_cv_v1", // 替换为你的索引名称
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
