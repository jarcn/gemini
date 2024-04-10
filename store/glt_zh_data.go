package store

import (
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type GltZhData struct {
	ID      string `json:"id" db:"id"`
	URL     string `json:"cv_url" db:"cv_url"`
	Profile string `json:"cv_proflie" db:"cv_proflie"`
}

func (cv *GltZhData) SelectAllData(db *sqlx.DB) ([]GltZhData, error) {
	querySql := `select * from qiyee_job_data.tbl_glt_zh_profile order by id asc limit 6000,2000`
	var result []GltZhData
	err := db.Select(&result, querySql)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cv *GltZhData) SelectByPage(start, end int, db *sqlx.DB) ([]GltZhData, error) {
	querySql := `select * from qiyee_job_data.tbl_glt_zh_profile order by id desc limit ?,?`
	var result []GltZhData
	err := db.Select(&result, querySql, start, end)
	if err != nil {
		return nil, err
	}
	return result, nil
}
