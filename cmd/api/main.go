package main

import (
	"log"
	"os"

	"github.com/jitendravee/clean_go/internals/db"
	"github.com/jitendravee/clean_go/internals/store"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		godotenv.Load()
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Fatal("MONGO_URL environment variable not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "user"
	}

	cfg := config{
		addr: addr,
		db: dbConfig{
			addr:   mongoURL,
			dbName: dbName,
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.dbName)
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount(db)
	log.Fatal(app.run(mux))
}
