package resources

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

func (r *Resources) insertUrls() error {
	urlsData := strings.Split(os.Getenv("GEN_TARGETS"), ",")

	statement, _ := r.Db.Prepare("INSERT INTO urls (url) VALUES (?)")
	urls := make([]string, 0)
	for _, v := range urlsData {
		urls = append(urls, v)
	}

	for _, v := range urls {
		fmt.Println("Going to insert " + v)
		statement.Exec(v)
	}

	return nil
}

func (r *Resources) InitDb() error {
	r.Db, _ = sql.Open("sqlite3", "./myDb.db")
	statement, _ := r.Db.Prepare("CREATE TABLE IF NOT EXISTS urls (id INTEGER PRIMARY KEY, url varchar UNIQUE)")
	statement.Exec()
	truncate, _ := r.Db.Prepare("DELETE FROM urls")
	truncate.Exec()

	//r.insertUrls()

	return nil
}

