package initDB

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Db *sql.DB

func init() {
	var err error
	//Db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/ginhello")
	Db, err = sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}
