package main

import (
	"encoding/json"
	"gemini/db"
	"github.com/IBM/sarama"
	"log"
	"testing"
)

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
	var data = make(map[string]string)
	data["url"] = "http://10.128.0.250/18ebc71cac2e84a3c9fd38bd3f680d67.pdf"
	data["profile"] = `{"code": 200, "msg": "success", "data": {"resume_msg": " ", "img_msg": " Data Pribadi Tmp, Tgl Lahir : Majalengka, 23/12/1988 Jenis Kelamin u2014 : Laki -Laki Agama : Islam Warga Negara: Indonesia Status : Menikah Alamat : Perum Pesona Agila Raya blok. C2 nO. 17 karang bahagia Sukatani CIKARANG - BEKASI No. NIK + 32160223 12880003 Telepon ; 085833110402 Email zallazis392@gmail.com ECOMMERCE Promosi Soces Budgeting Go66 Komunikasi SoooG Komputerisasi MS Word oooos MS Exel Sooo MS Power Point S606 Internet oo66. DesignGraphic GOOG u2014 Kecakapan Bhs Indonesia ooooG Bhs Inggris 900 u2014 CURRICULUM VITAE Nama saya Abdul Azis, saya berumur 33 tahun. Saya berasal dari Majalengka, Jawa Barat, Pendidikan terakhir saya Diploma Ill di Politeknik Tri Mitra Karya Mandiri (TMKM) Karawang dengan jurusan tehnik komputer, saya lulus pada tanggal 18 November 2017. Sebelum dan Setelah saya lulus saya sudah mempunyai pengalaman bekerja di beberapa perusahaan. Pengalaman Kerja General Contraktor Sebuah perusahaan yang bergerak dalam bidang jasa Kontraktor, yang beralamat di margasari, Karawaci - Tangerang u00bb Sebagai Kepala LOGISTIK PT MITRA MAKMUR GEMILANG (2008 sd 2011 ) u00bb Operator Produksi ~ Kepala Grup Divisi PT ADYAWINSA DINAMIKA KARAWANG (2011 sd 2013 ) u00bb Assembly PT. ICHIKOH INDONESIA (2013 sd 2014) u00bb Wharehouse u00bb Delivery 1.0 EVENT ORGANIZER ( 2014 sd 2015 ) 7 Pengarah Acara dan Desing Acara PT. BINAKARYA GROUP ( 4 bulan- 2016) >u00bb Marketing Eksekutif RM, ANEKA SAMBAL ( 2016 sd 2017 ) > Team Kreatif > Supervisior PT. PRISINDO PRIMA UTAMA (3 bulan -2018 ) > Staff ITC (ecommerce) LIPPO GROUP ( 4 bulan 2018 } >u00bb Marketing Eksekutif PT SICEPAT EKPRESS ( 5 Tahun 2018 sd 2023) > Kepala Cabang > Sigesit Pickup | Pendidikan FORMAL Politeknik Trimitra Karya Mandiri u00bb Diploma 3 Jurusan Tehnik Informatika 2017 y IPK:3.14 SMK BINASWASTA KUNINGAN u00bb PENJUALAN Tahun 2006 Dipindal dengan CamSeanner f", "id": "7402c8751fcdbc57561f8a647b04eef9"}}`
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
