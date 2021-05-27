package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"todoapi/handler"
	"todoapi/infra/persistence"
	"todoapi/usecase"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func main() {
	// server config
	addr := "0.0.0.0:4000"
	server := &http.Server{
		Addr: addr,
	}

	log.Print("Connecting db...")
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("DB ready!!")

	itemRepository := persistence.NewItemSqlRepository(db)
	itemUseCase := usecase.NewItemUseCase(itemRepository)
	itemHandler := handler.NewItemHandler(itemUseCase)

	// add handlers
	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/items", itemHandler.HandleAll)
	http.HandleFunc("/items/", itemHandler.HandleOne)

	// start server
	server.ListenAndServe()
}

func initDB() (*sql.DB, error) {
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
