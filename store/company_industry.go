package store

import (
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type CompanyIndustry struct {
	ID              string `json:"id" db:"id"`
	CompanyName     string `json:"company_name" db:"company_name"`
	CompanyIndustry string `json:"company_industry" db:"company_industry"`
}

func (cv *CompanyIndustry) SelectAllData(db *sqlx.DB) ([]CompanyIndustry, error) {
	querySql := `select * from qiyee_job_data.tbl_company_industry`
	var result []CompanyIndustry
	err := db.Select(&result, querySql)
	if err != nil {
		return nil, err
	}
	return result, nil
}
