package db

import (
	"github.com/cookieY/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
	log.Println("mysql init success")
}

func Client() *sqlx.DB {
	return mdb
}

func Close() {
	mdb.Close()
}
