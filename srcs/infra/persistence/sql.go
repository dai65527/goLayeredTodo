package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func InitDB() (*sql.DB, error) {
	var dbSource string
	var dbDriver string

	// create db connection
	if os.Getenv("DB_MIDDLEWARE") == "mysql" {
		dbDriver = "mysql"
		dbName := os.Getenv("MYSQL_DATABASE")
		dbUser := os.Getenv("MYSQL_USER")
		dbPass := os.Getenv("MYSQL_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost:3306"
		}
		dbSource = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)
	} else {
		dbDriver = "sqlite"
		dbSource = "./database.db"
	}
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	// wait for db start
	count := 30
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
		if count < 1 {
			return nil, err
		}
		count--
	}

	// init database
	var sql string
	if os.Getenv("DB_MIDDLEWARE") == "mysql" {
		sql = `
			CREATE TABLE IF NOT EXISTS items (
				id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
				name TEXT NOT NULL,
				done BOOLEAN NOT NULL DEFAULT 0
			);`
	} else {
		sql = `
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0
		);`
	}
	_, err = db.Exec(sql)
	if err != nil {
		return nil, err
	}
	return db, err
}
