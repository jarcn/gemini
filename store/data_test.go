package store

import (
	"fmt"
	"gemini/db"
	"testing"
)

func TestQuery(t *testing.T) {
	db.MustInitMySQL(db.MysqlAddr)
	cvData := CvData{}
	cv, err := cvData.GetCvByUrl(db.Client(), "http://10.128.0.250/18ebc71cac2e84a3c9fd38bd3f680d67.pdf")
	if err != nil {
		fmt.Println(err)
	}
	msg := cv.ResumeMsg
	fmt.Println(msg)
}

func TestQueryById(t *testing.T) {
	db.MustInitMySQL(db.MysqlAddr)
	cvData := GeminiResult{}
	cv, err := cvData.QueryById(db.Client(), 1)
	if err != nil {
		fmt.Println(err)
	}
	msg := cv.Type
	fmt.Println(msg)
}
