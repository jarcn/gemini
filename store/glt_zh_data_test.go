package store

import (
	"encoding/json"
	"gemini/db"
	"github.com/IBM/sarama"
	"log"
	"testing"
)

func init() {
	db.MustInitMySQL(db.MysqlAddr)
}

func TestProd(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(db.KafkaBrokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()
	var d = GltZhData{}
	allData, _ := d.SelectAllData(db.Client())
	for _, datum := range allData {
		data := make(map[string]string)
		data["url"] = datum.URL
		data["profile"] = datum.Profile
		marshal, _ := json.Marshal(data)
		msg := &sarama.ProducerMessage{ // 构建消息
			Topic: db.Step1Topic,
			Value: sarama.StringEncoder(marshal),
		}
		log.Println(msg)
		partition, offset, err := producer.SendMessage(msg) // 发送消息
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
		log.Printf("Message sent successfully! Partition: %d, Offset: %d", partition, offset)
	}
}
