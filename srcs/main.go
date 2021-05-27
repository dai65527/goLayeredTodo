package main

import (
	"log"
	"net/http"
	"os"
	"todoapi/domain/repository"
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

	var itemRepository repository.ItemRepository
	if os.Getenv("DB_MONGO") == "on" {
		log.Println("connecting to mongoDB...")
		db, err := persistence.InitMongo()
		if err != nil {
			log.Fatal(err)
		}
		itemRepository = persistence.NewItemMongoRepository(db)
	} else {
		db, err := persistence.InitDB()
		if err != nil {
			log.Fatal(err)
		}
		itemRepository = persistence.NewItemSqlRepository(db)
	}
	log.Print("DB ready!!")

	itemUseCase := usecase.NewItemUseCase(itemRepository)
	itemHandler := handler.NewItemHandler(itemUseCase)

	// add handlers
	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/items", itemHandler.HandleAll)
	http.HandleFunc("/items/", itemHandler.HandleOne)

	// start server
	server.ListenAndServe()
}
