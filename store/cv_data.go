package store

import (
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type CvData struct {
	ID         string    `json:"id" db:"id"`
	URL        string    `json:"url" db:"url"`
	ResumeMsg  string    `json:"resume_msg" db:"resume_msg"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

func (cv *CvData) GetCvByUrl(db *sqlx.DB, url string) (*CvData, error) {
	querySql := ` select * from qiyee_job_data.tbl_resume_parse_data where url = ?`
	var result []CvData
	err := db.Select(&result, querySql, url)
	if err != nil {
		return nil, err
	}
	if len(result) >= 1 {
		return &result[0], nil
	} else {
		return nil, nil
	}
}
