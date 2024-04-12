package tasks

import (
	"gemini/db"
	"gemini/store"
	"testing"
)

func init() {
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //生产环境
}

func TestDoMergeBeat(t *testing.T) {
	var d = store.GeminiResult{}
	all, _ := d.FindAllError(db.Client())
	for _, result := range all {
		ReDoMergeBeat(*result, "111111")
	}

}
