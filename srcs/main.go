package main

import (
	"log"
	"net/http"
	"os"
	"todoapi/handler"
	"todoapi/infra/persistence"
	"todoapi/usecase"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func main() {
	// server config
	addr := os.Getenv("API_HOST")
	if os.Getenv("API_HOST") == "" {
		addr = "0.0.0.0:4000"
	}
	server := &http.Server{
		Addr: addr,
	}

	log.Print("Connecting db...")
	db, err := persistence.InitDB()
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
