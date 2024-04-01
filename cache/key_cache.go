package cache

import (
	"fmt"
	"gemini/db"
	"gemini/store"
	"math/rand"
	"time"
)

var KeyCache []string

func InitKeyCache() {
	var keyArr []KeyInfo
	query := `select api_key from qiyee_job_data.tbl_gemini_keys;`
	db.Client().Select(&keyArr, query)
	KeyCache = make([]string, len(keyArr))
	for i := 0; i < len(keyArr); i++ {
		KeyCache[i] = keyArr[i].Key
	}
	fmt.Println("gemini key init success")
}

func GetKey() string {
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(KeyCache))
	key := KeyCache[randomIndex]
	result := store.GeminiResult{}
	count, err := result.CountByKey(db.Client(), key)
	if err != nil {
		return GetKey()
	}
	if count == 0 {
		return key
	}
	currentTime := time.Now().Unix()
	if currentTime-count > 60 {
		return key
	} else {
		time.Sleep(time.Second * 60)
	}
	return GetKey()
}

type KeyInfo struct {
	Key string `db:"api_key" json:"key"`
}
