package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS emails(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		email TEXT,
		msgId INTEGER,
		UNIQUE(email,msgId)
	  );`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO emails (email, msgId) values ($1, $2)`,
		"kiling@mail.ru", 8035)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT COUNT(*) FROM emails WHERE email=$1 AND msgId=$2",
		"kiling@mail.ru", 8035)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(count)
}
