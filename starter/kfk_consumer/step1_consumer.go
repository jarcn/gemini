package kfk_consumer

import (
	"context"
	"gemini/cache"
	"gemini/tasks"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
)

func Step1ConsumerStart(brokers []string, topic, group string) {
	config := sarama.NewConfig()                                           // 创建 Kafka 消费者配置
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange // 使用Range重新平衡策略
	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)  // 创建消费者组
	if err != nil {
		log.Println("Error creating consumer group:", err)
		return
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Println("Error closing consumer group:", err)
		}
	}()
	handler := &Step1ConsumerGroupHandler{}                 // 处理函数
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
	log.Println("Step1Consumer Starting ...")
	sigterm := make(chan os.Signal, 1) // 等待中断信号以优雅地关闭消费者
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm
	log.Println("Interrupt received, shutting down consumer group...")
	cancel()
	wg.Wait()
	log.Println("Consumer group shutdown complete.")
}

// ConsumerGroupHandler 实现了sarama.ConsumerGroupHandler接口
type Step1ConsumerGroupHandler struct{}

// Setup 在消费者组开始消费前调用，用于初始化
func (Step1ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup 在消费者组停止消费后调用，用于清理
func (Step1ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 将消息传递给处理程序处理
func (h *Step1ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if h.ProcessMessage(message) { //消息处理
			session.MarkMessage(message, "gemini parse success") // 提交消息的消费确认信息
		} else {
			log.Println("Message processing failed.") // 处理失败，不提交确认信息
		}
	}
	return nil
}

func (h *Step1ConsumerGroupHandler) ProcessMessage(message *sarama.ConsumerMessage) bool {
	log.Printf("Message claimed: partition = %d, offset = %d, topic = %s\n", message.Partition, message.Offset, message.Topic)
	return tasks.DoMerge(message.Value, cache.GetKey())
}
