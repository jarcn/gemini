package main

import (
	"encoding/json"
	"gemini/db"
	"github.com/IBM/sarama"
	"log"
	"testing"
)

func TestProd(t *testing.T) {
	product()
}

func product() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(db.KafkaBrokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()
	var data = make(map[string]string)
	data["url"] = "http://10.128.0.250/18ebc71cac2e84a3c9fd38bd3f680d67.pdf"
	data["profile"] = "{\"id\":\"0df7eeb1-f79d-45d5-af76-36e96bb46ff6\",\"talentProfileSummary\":{\"id\":\"4ddd5d62-0904-49f8-a47f-f2847147a459\",\"MergedTalentProfileId\":\"0df7eeb1-f79d-45d5-af76-36e96bb46ff6\",\"sources\":[\"GLINTS\"],\"resume\":\"18ebc71cac2e84a3c9fd38bd3f680d67.pdf\",\"experiences\":[{\"title\":\"Warehouse\",\"organization\":\"PT. ICHIKOH indonesia\",\"description\":null,\"startDate\":\"2013-05-01\",\"endDate\":\"2014-09-30\"}],\"skills\":[{\"id\":\"41c6aabc-1aa8-4180-bc91-39bd91eb1a04\",\"name\":\"\\tJujur, teliti, dan rajin\"},{\"id\":\"44cdd62b-8ee9-480c-8f12-feb2a502bd8c\",\"name\":\"Memiliki Kemampuan Komunikasi Yang Baik\"},{\"id\":\"7b134fed-13ab-4ed4-8acd-918c9599379e\",\"name\":\"Mampu Bekerja Keras Dan Jujur Serta Cepat Mampu Menyesuaikan Dan\"},{\"id\":\"9934ce96-8bba-4bfa-9beb-78acdb971388\",\"name\":\"\\tBisa mengoperasikan komputer Wi\"},{\"id\":\"b5dd5f98-2876-4efe-8e3e-f534d2a24f0c\",\"name\":\"Saya bisa berkerja secara individu maupun tim menyukai tantangan dan bisa berkerja di bawah tekanan\"}],\"name\":\"Abdul Azis\",\"email\":\"allazis392@gmail.com\",\"location\":\"Bekasi, Indonesia\",\"lastSeen\":\"2023-10-22T13:05:22.203Z\",\"birthDate\":\"\",\"phoneNumbers\":[\"+6285833110402\"],\"nationality\":\"\",\"salary\":{\"currencyCode\":\"IDR\",\"latest\":null,\"expectation\":4000000},\"candidateStatus\":\"I_AM_LOOKING_FOR_JOB\",\"updatedAt\":\"2023-10-22T13:05:22.203Z\",\"profilePics\":[{\"key\":null,\"source\":\"GLINTS\"}],\"network\":{\"website\":\"\",\"LinkedIn\":\"\",\"CakeResume\":\"\",\"Facebook\":\"\",\"Twitter\":\"\",\"Instagram\":\"\",\"Behance\":\"\",\"GitHub\":\"\",\"CodePen\":\"\",\"Vimeo\":\"\",\"Youtube\":\"\"},\"intro\":null,\"educations\":[{\"degree\":\"Amd/d3\",\"school\":\"Politeknik Trimitra Karya Mandiri\",\"description\":null,\"fieldOfStudy\":\"Tehnik Informatika\",\"startDate\":\"2015-11-01\",\"endDate\":\"2017-11-30\"}]},\"saved\":null}"
	marshal, _ := json.Marshal(data)
	msg := &sarama.ProducerMessage{ // 构建消息
		Topic: db.Step1Topic,
		Value: sarama.StringEncoder(marshal),
	}
	partition, offset, err := producer.SendMessage(msg) // 发送消息
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Message sent successfully! Partition: %d, Offset: %d", partition, offset)
}
