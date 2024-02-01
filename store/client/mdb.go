package client

import (
	_ "github.com/go-sql-driver/mysql"
	"log"

	"github.com/cookieY/sqlx"
)

var mdb *sqlx.DB

func MustInitMySQL(datasource string) {
	db, err := sqlx.Open("mysql", datasource)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	mdb = db
	log.Println("mysql init success ...")
}

func DbClient() *sqlx.DB {
	return mdb
}

func Close() {
	mdb.Close()
}
