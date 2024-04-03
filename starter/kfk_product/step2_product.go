package kfk_product

import (
	"github.com/IBM/sarama"
	"log"
)

func SendMsg(brokers []string, topic string, data []byte) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{ // 构建消息
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	partition, offset, err := producer.SendMessage(msg) // 发送消息
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Message sent successfully! Partition: %d, Offset: %d", partition, offset)
}
