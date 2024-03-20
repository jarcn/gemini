package cache

import (
	"fmt"
	"gemini/db"
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
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(KeyCache))
	return KeyCache[randomIndex]
}

type KeyInfo struct {
	Key string `db:"api_key" json:"key"`
}
