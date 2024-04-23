package store

import (
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type HraCvData struct {
	ResumeLink string `json:"resume_link" db:"resume_link"`
	Content    string `json:"content" db:"content"`
}

func (cv *HraCvData) SelectAll(db *sqlx.DB) ([]HraCvData, error) {
	querySql := `select distinct a.resume_link,b.content from crmdb.tbl_hr_upload_cv_record a join crmdb.tbl_order_external_js_cv_content b on a.external_js_id = b.external_js_id and a.create_time > "2024-03-26 10:11:53";`
	var result []HraCvData
	err := db.Select(&result, querySql)
	if err != nil {
		return nil, err
	}
	if len(result) >= 1 {
		return result, nil
	} else {
		return nil, nil
	}
}
