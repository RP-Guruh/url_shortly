package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Data struct {
	original_url string
	short_url    string
	unique_key   string
	created_at   time.Time
}

func connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db/shorturl_db.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func saveurltodb(originalurl string, shorturl string, uniquekey string) {

	db, err := connection()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	defer db.Close()

	data := &Data{
		original_url: originalurl,
		short_url:    shorturl,
		unique_key:   uniquekey,
		created_at:   time.Now(),
	}

	_, err = db.Exec(
		"INSERT INTO yourlink (original_url, short_url, unique_key, created_at) VALUES (?, ?, ?, ?)",
		data.original_url,
		data.short_url,
		data.unique_key,
		data.created_at,
	)

	if err != nil {
		log.Fatal("Failed to prepare insert query:", err)
		fmt.Println("Failed to prepare insert query:", err)
		return
	}
	log.Println("Data saved to the database successfully")
	fmt.Println("Data saved to the database successfully")

}

func getOriginalURL(uniqueKey string) (string, error) {
	db, err := connection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return "", err
	}
	defer db.Close()

	var originalURL string
	err = db.QueryRow("SELECT original_url FROM yourlink WHERE unique_key = ?", uniqueKey).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		log.Printf("Failed to query the database: %v", err)
		return "", err
	}

	return originalURL, nil
}
