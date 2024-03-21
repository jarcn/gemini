package main

import (
	"encoding/json"
	"fmt"
	"gemini/cache"
	"gemini/db"
	"gemini/store"
	"gemini/tasks"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
	"sync"
)

func init() {
	//db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //预发环境
	cache.InitKeyCache()
}

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	//consumer, err := sarama.NewConsumer([]string{"10.128.0.94:9092", "10.128.0.156:9092", "10.128.0.124:9092"}, config) //生产环境
	consumer, err := sarama.NewConsumer([]string{"10.129.0.78:9092", "10.129.0.180:9092", "10.129.0.85:9092"}, config) //预发环境
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println("Error closing consumer:", err)
		}
	}()
	// 指定要消费的主题
	topic := "go-profile-merge"
	// 指定要消费的分区，这里为空表示消费所有分区
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error retrieving partition list:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(partitionList))
	// 遍历每个分区创建消费者
	for _, partition := range partitionList {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("Error creating partition consumer for partition %d: %v", partition, err)
			continue
		}
		// 异步处理每个分区的消息
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for {
				select {
				case msg := <-pc.Messages():
					doMerge(msg.Value)
				case err := <-pc.Errors():
					fmt.Println("Error:", err)
				}
			}
		}(partitionConsumer)
	}
	// 等待信号以优雅地关闭消费者
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm
	fmt.Println("Interrupt received, shutting down consumer...")
	wg.Wait()
	fmt.Println("Consumer shutdown complete.")
}

func doMerge(msg []byte) {
	data := make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	key := cache.GetKey()
	profile, _ := data["profile"].(string)
	url, _ := data["url"].(string)
	result := store.GeminiResult{
		GeminiKey:   key,
		ProfileData: profile,
		CVURL:       url,
		CVData:      profile,
	}
	id, err := result.Create(db.Client())
	if err != nil {
		fmt.Println("insert data error:", err)
		return
	}
	step1 := tasks.CallGemini(profile, profile, key)
	result.GeminiStep1 = step1
	result.ID = id
	result.Update(db.Client())
}
