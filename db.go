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

func storeVisitor(id_url string, visitorInfo VisitorInfo) (string, error) {
	// Buka koneksi database
	db, err := connection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return "", err
	}
	defer db.Close()

	// Query untuk insert data, termasuk created_at
	insertQuery := `
		INSERT INTO visitorrow (
			query_url, ip_address, browser, device, os, city, region, 
			country, country_name, continent, continent_name, latlong, 
			org, postal, timezone, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'));
	`

	// Eksekusi insert
	_, err = db.Exec(insertQuery,
		id_url,
		visitorInfo.IPAddress,
		visitorInfo.Browser,
		visitorInfo.Device,
		visitorInfo.OS,
		visitorInfo.City,
		visitorInfo.Region,
		visitorInfo.Country,
		visitorInfo.CountryName,
		visitorInfo.Continent,
		visitorInfo.ContinentName,
		visitorInfo.LatLong,
		visitorInfo.Org,
		visitorInfo.Postal,
		visitorInfo.Timezone,
	)

	if err != nil {
		log.Printf("Failed to execute insert query: %v", err)
		return "", err
	}

	log.Printf("Data saved to the database successfully for ID: %s", id_url)
	return id_url, nil
}

func getMyVisitor(uniqueKey string) (int, error) {
	db, err := connection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return 0, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM visitorrow WHERE query_url = ?", uniqueKey).Scan(&count)
	if err != nil {
		log.Println("Error saat menghitung data:", err)
		return 0, err
	}

	return count, nil
}
