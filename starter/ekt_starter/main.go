package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gemini/profile"
	"gemini/utils"
	"github.com/buger/jsonparser"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	root := "/Users/tarzan/Documents/ekt"
	var dataArr []profile.Talent
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
			data := parseData(content)
			dataArr = append(dataArr, data...)
			if len(dataArr) == 1000 {
				save2Es(dataArr)
				dataArr = []profile.Talent{}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", root, err)
	}
}

func parseData(content []byte) []profile.Talent {
	var talents []profile.Talent
	jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var data profile.Talent
		json.Unmarshal(value, &data)
		formatData(&data)
		talents = append(talents, data)
	}, "data")
	return talents
}

func formatData(data *profile.Talent) {
	data.BirthDate = utils.DateFormat(data.BirthDate)
	educations := data.Educations
	for i := 0; i < len(educations); i++ {
		education := &educations[i]
		education.StartDate = utils.DateFormat(education.StartDate)
		education.EndDate = utils.DateFormat(education.EndDate)
	}
	experiences := data.Experiences
	for i := 0; i < len(experiences); i++ {
		experience := &experiences[i]
		experience.StartDate = utils.DateFormat(experience.StartDate)
		experience.EndDate = utils.DateFormat(experience.EndDate)
		experience.Description = removeHTMLTags(experience.Description)
	}
}

func removeHTMLTags(text string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(text, "")
}

func save2Es(talents []profile.Talent) {
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
	for _, doc := range talents {
		// 将文档转换为 JSON 格式
		docJSON, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling document: %s", err)
			continue
		}

		// 准备批量插入请求
		req := esapi.IndexRequest{
			Index:      "ekt_data", // 替换为你的索引名称
			DocumentID: strconv.Itoa(doc.Id),
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
			log.Printf("[%s] Indexed document with ID: %d", res.Status(), doc.Id)
		}
	}
}
