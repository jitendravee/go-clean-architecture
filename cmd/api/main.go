package main

import (
	"log"
	"os"

	"github.com/jitendravee/clean_go/internals/db"
	"github.com/jitendravee/clean_go/internals/store"
)

func main() {
	// Fetch environment variables
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080" // Default port
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Fatal("MONGO_URL environment variable not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "user" // Default database name
	}

	// Set up configuration and initialize the database connection
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

	// Set up application and routes
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount(db)

	// Start the server
	log.Fatal(app.run(mux))
}
