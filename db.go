package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var gAskDB *sql.DB

// GetAskDB 获取当前的数据库连接
func GetAskDB() *sql.DB {
	if gAskDB == nil {
		db, err := sql.Open("sqlite3", "./ask.db")
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(`
create table if not exists cards (
	id char(50) primary key not null,
	title char(50) not null,
	question text not null,
	answer text,
	review_time datetime default current_timestamp,
	level integer default 0,
	hash char(16)
)`)
		if err != nil {
			log.Fatal(err)
		}
		gAskDB = db
	}
	return gAskDB
}
