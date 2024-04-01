package store

import (
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type GeminiResult struct {
	ID          int64  `json:"id" db:"id"`
	CVData      string `json:"cv_data" db:"cv_data"`                         //原始pdf中ocr识别的数据
	ProfileData string `json:"profile_data" db:"profile_data"`               //原始profile结构化的数据
	CVURL       string `json:"cv_url" db:"cv_url"`                           //CV pdf 存储地址
	GeminiStep1 string `json:"gemini_step1_result" db:"gemini_step1_result"` //任务一结构化的数据
	GeminiStep2 string `json:"gemini_step2_result" db:"gemini_step2_result"` //任务二推导的结果
	GeminiStep3 string `json:"gemini_step3_result" db:"gemini_step3_result"` //任务三推导的结果
	GeminiStep4 string `json:"gemini_step4_result" db:"gemini_step4_result"` //存储step1和step2的合并结果
	CreateTime  int64  `json:"create_time" db:"create_time"`
	UpdateTime  int64  `json:"update_time" db:"update_time"`
	GeminiKey   string `json:"gemini_key" db:"gemini_key"`
}

func (gr *GeminiResult) Create(db *sqlx.DB) (int64, error) {
	currentTime := time.Now().Unix()
	insertQuery := `INSERT INTO tbl_gemini_result (cv_data, profile_data, cv_url, gemini_step1_result, gemini_step2_result, gemini_step3_result, gemini_step4_result, create_time, update_time, gemini_key)
        			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(insertQuery, gr.CVData, gr.ProfileData, gr.CVURL, gr.GeminiStep1, gr.GeminiStep2, gr.GeminiStep3, gr.GeminiStep4, currentTime, currentTime, gr.GeminiKey)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (gr *GeminiResult) Read(db *sqlx.DB, id int64) error {
	readQuery := ` SELECT * FROM tbl_gemini_result WHERE id=?`
	return db.QueryRow(readQuery, id).Scan(
		&gr.ID, &gr.CVData, &gr.ProfileData, &gr.CVURL, &gr.GeminiStep1, &gr.GeminiStep2, &gr.GeminiStep3, &gr.GeminiStep4, &gr.CreateTime, &gr.UpdateTime, &gr.GeminiKey,
	)
}

func (gr *GeminiResult) CountByKey(db *sqlx.DB, key string) (int64, error) {
	countSql := ` select * from qiyee_job_data.tbl_gemini_result where gemini_key = ? order by create_time desc limit 1`
	var result []GeminiResult
	err := db.Select(&result, countSql, key)
	if err != nil {
		return -1, err
	}
	if len(result) >= 1 {
		return result[0].CreateTime, nil
	} else {
		return 0, nil
	}
}

func (gr *GeminiResult) CvExists(db *sqlx.DB, cvUrl string) (bool, error) {
	countSql := ` select * from qiyee_job_data.tbl_gemini_result where cv_url = ?`
	var result []GeminiResult
	err := db.Select(&result, countSql, cvUrl)
	if err != nil {
		return false, err
	}
	if len(result) >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (gr *GeminiResult) QueryById(db *sqlx.DB, id int64) (*GeminiResult, error) {
	countSql := ` select * from qiyee_job_data.tbl_gemini_result where id = ?`
	var result []GeminiResult
	err := db.Select(&result, countSql, id)
	if err != nil {
		return nil, err
	}
	if len(result) >= 1 {
		return &result[0], nil
	} else {
		return nil, nil
	}
}

func (gr *GeminiResult) Update(db *sqlx.DB) error {
	currentTime := time.Now().Unix()
	updateQuery := `UPDATE tbl_gemini_result
        			SET cv_data=?, profile_data=?, cv_url=?, gemini_step1_result=?, gemini_step2_result=?, gemini_step3_result=?, gemini_step4_result=?, update_time=?, gemini_key=?
        			WHERE id=?`
	_, err := db.Exec(updateQuery, gr.CVData, gr.ProfileData, gr.CVURL, gr.GeminiStep1, gr.GeminiStep2, gr.GeminiStep3, gr.GeminiStep4, currentTime, gr.GeminiKey, gr.ID)
	return err
}

func (gr *GeminiResult) Step2Update(db *sqlx.DB) error {
	currentTime := time.Now().Unix()
	updateQuery := `UPDATE tbl_gemini_result SET gemini_step2_result=?, update_time=?, gemini_key=? WHERE id=?`
	_, err := db.Exec(updateQuery, gr.GeminiStep2, currentTime, gr.GeminiKey, gr.ID)
	return err
}

func (gr *GeminiResult) Step4Update(db *sqlx.DB) error {
	currentTime := time.Now().Unix()
	updateQuery := `UPDATE tbl_gemini_result SET gemini_step4_result=?, update_time=?, gemini_key=? WHERE id=?`
	_, err := db.Exec(updateQuery, gr.GeminiStep4, currentTime, gr.GeminiKey, gr.ID)
	return err
}

func (gr *GeminiResult) Delete(db *sqlx.DB) error {
	deleteQuery := ` DELETE FROM tbl_gemini_result WHERE id=? `
	_, err := db.Exec(deleteQuery, gr.ID)
	return err
}

func (gr *GeminiResult) FindAll(db *sqlx.DB) ([]GeminiResult, error) {
	findAll := `select distinct(cv_url),gemini_step1_result,id from qiyee_job_data.tbl_gemini_result where gemini_step1_result !=""`
	var result []GeminiResult
	err := db.Select(&result, findAll)
	return result, err
}
