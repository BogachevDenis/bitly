package database

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)
var db *sql.DB

func Connect (pg_user, pg_pass, pg_base string ) error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", pg_user, pg_pass, pg_base))
	if err != nil {
		log.WithFields(log.Fields{
			"Connect" : "database",
			"error" : err,
		}).Fatal("Connection error")
		return err
	}
	return nil
}

func GetlongURL(shorturl string) string {
	r := db.QueryRow(`SELECT "longurl" FROM "url" WHERE "shorturl" = $1`, shorturl)
	var longurl string
	err := r.Scan(&longurl)
	if err == sql.ErrNoRows {
		return ""
	}
	return longurl
}

func InsertData(longURL, shortURL string) error  {
	_, error := db.Exec(`INSERT INTO url (longurl, shorturl) VALUES($1,$2)`,longURL, shortURL)
	if error != nil {
		log.WithFields(log.Fields{
			"error" : error,
		}).Fatal("Insert error")
		return error
	}
	return nil
}