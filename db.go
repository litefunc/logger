package logger

import (
	"database/sql"
	"log"
)

func SetDb(db *sql.DB) {
	defaultLogger.DB = db
}

func saveToDb(db *sql.DB, logs logs) error {
	if _, err := db.Exec(`insert into log(time, level, service, file, line, msg) values($1, $2, $3, $4, $5, $6)`, logs.Ltime, logs.Level, logs.Service, logs.Lfile, logs.Lline, logs.Msg); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
