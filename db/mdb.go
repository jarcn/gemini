package db

import (
	"fmt"
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

var mdb *sqlx.DB

func MustInitMySQL(dbUrl string) {
	db, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	mdb = db
	fmt.Println("mysql init success")
}

func Client() *sqlx.DB {
	return mdb
}

func Close() {
	mdb.Close()
}
