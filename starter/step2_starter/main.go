package main

import (
	"context"
	"encoding/json"
	"gemini/cache"
	"gemini/db"
	"gemini/store"
	"gemini/tasks"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
)

func init() {
	db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
	//db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //预发环境
	cache.InitKeyCache()
}

var brokers = []string{"10.128.0.94:9092", "10.128.0.156:9092", "10.128.0.124:9092"} // 生产环境
// var brokers = []string{"10.128.0.94:9092", "10.128.0.156:9092", "10.128.0.124:9092"} // 预发环境
var topic = "gemini-step2" // 要消费的主题
var consumerGroup = "gemini-step2-group"

func main() {
	//productStart()
	//consumerStart()
	consumerSync()
}

func consumerSync() {
	var result store.GeminiResult
	allData, _ := result.FindByIds([]int64{3, 5, 6}, db.Client())
	for _, d := range allData {
		var data = make(map[string]int64)
		data["id"] = d.ID
		marshal, _ := json.Marshal(data)
		tasks.DoDeduce(marshal, "AIzaSyBTRLUbz_9wX_prQWNjjtLueMIaF_uJEm8", false)
	}
}

// productStart 先写数据(自产自销)
func productStart() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()
	var result store.GeminiResult
	allData, err := result.SelectAll(db.Client())
	for _, d := range allData {
		var data = make(map[string]int64)
		data["id"] = d.ID
		marshal, _ := json.Marshal(data)
		msg := &sarama.ProducerMessage{ // 构建消息
			Topic: topic,
			Value: sarama.StringEncoder(marshal),
		}
		go producer.SendMessage(msg) // 发送消息
		//if err != nil {
		//	log.Fatalf("Failed to send message: %v", err)
		//}
		//log.Printf("Message sent successfully! Partition: %d, Offset: %d", partition, offset)
	}
}

// consumerStart 启动消费者
func consumerStart() {
	config := sarama.NewConfig()                                                  // 创建 Kafka 消费者配置
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange        // 使用Range重新平衡策略
	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerGroup, config) // 创建消费者组
	if err != nil {
		log.Println("Error creating consumer group:", err)
		return
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Println("Error closing consumer group:", err)
		}
	}()
	handler := &ConsumerGroupHandler{}                      // 处理函数
	ctx, cancel := context.WithCancel(context.Background()) // 开始消费消息
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				log.Println("Error from consumer:", err)
				cancel()
				return
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	// 等待中断信号以优雅地关闭消费者
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm
	log.Println("Interrupt received, shutting down consumer group...")
	cancel()
	wg.Wait()
	log.Println("Consumer group shutdown complete.")
}

// 实现了sarama.ConsumerGroupHandler接口
type ConsumerGroupHandler struct{}

// 在消费者组开始消费前调用，用于初始化
func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// 在消费者组停止消费后调用，用于清理
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// 将消息传递给处理程序处理
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if h.ProcessMessage(message) { //消息处理
			session.MarkMessage(message, "gemini parse success") // 提交消息的消费确认信息
		} else {
			log.Println("Message processing failed.") // 处理失败，不提交确认信息
		}
	}
	return nil
}

func (h *ConsumerGroupHandler) ProcessMessage(message *sarama.ConsumerMessage) bool {
	return true
	//return tasks.DoDeduce(message.Value, "AIzaSyBTRLUbz_9wX_prQWNjjtLueMIaF_uJEm8", false)
}
