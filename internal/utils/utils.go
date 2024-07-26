package utils

import (
	"database/sql"
	"log"
	"os"
)

func ReadStringFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func GetWeatherEmoji(temp int) string {
	switch {
	case temp <= 13:
		return "❄️"
	case temp > 13 && temp <= 20:
		return "🧥"
	case temp > 20 && temp <= 30:
		return "👕"
	case temp > 30:
		return "💀"
	default:
		return "???"
	}
}

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../storage/streamers.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS streamers (
			twitch TEXT PRIMARY KEY,
			oauth_token TEXT,
			gapi TEXT,
			mapi TEXT
		)
	`

	_, err := db.Exec(stmt)
	return err
}

func Contains(db *sql.DB, twitch string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM streamers WHERE twitch = ? LIMIT 1)`
	err := db.QueryRow(query, twitch).Scan(&exists)
	return exists, err
}

func Insert(db *sql.DB, twitch, oauthToken, gapi, mapi string) error {
	stmt, err := db.Prepare(`INSERT INTO streamers (twitch, oauth_token, gapi, mapi) 
	VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(twitch, oauthToken, gapi, mapi)
	return err
}
