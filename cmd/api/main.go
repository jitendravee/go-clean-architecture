package main

import (
	"log"
	"os"

	"github.com/jitendravee/clean_go/internals/db"
	"github.com/jitendravee/clean_go/internals/store"
	"github.com/rs/cors"
)

// func main() {

// 	envPath := flag.String("env", ".env", "Path to .env file")
// 	flag.Parse()

// 	// Load .env file
// 	err := godotenv.Load(*envPath)
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Get environment variables
// 	addr := os.Getenv("ADDR")
// 	if addr == "" {
// 		addr = ":8080" // Default port
// 	}

// 	mongoURL := os.Getenv("MONGO_URL")
// 	if mongoURL == "" {
// 		log.Fatal("MONGO_URL environment variable not set")
// 	}

// 	dbName := os.Getenv("DB_NAME")
// 	if dbName == "" {
// 		dbName = "user" // Default database name
// 	}

// 	// Set up configuration
// 	cfg := config{
// 		addr: addr,
// 		db: dbConfig{
// 			addr:   mongoURL,
// 			dbName: dbName,
// 		},
// 	}

// 	// Initialize database connection
// 	db, err := db.New(cfg.db.addr, cfg.db.dbName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	store := store.NewStorage(db)

// 	// Set up application and routes
// 	app := &application{
// 		config: cfg,
// 		store:  store,
// 	}
// 	mux := app.mount(db)

// 	// Enable CORS with middleware
// 	handler := cors.Default().Handler(mux) // Default CORS settings allow all origins

//		// Start the server with CORS enabled
//		log.Fatal(app.run(handler))
//	}
func main() {
	// Get environment variables directly from the OS
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080" // Default port
	}

	mongoURL := "mongodb+srv://jitendrajat6397:Jitendra6323@jitendra.yaofk.mongodb.net/"
	if mongoURL == "" {
		log.Fatal("MONGO_URL environment variable not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "user" // Default database name
	}

	// Set up configuration
	cfg := config{
		addr: addr,
		db: dbConfig{
			addr:   mongoURL,
			dbName: dbName,
		},
	}

	// Initialize database connection
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

	// Enable CORS with middleware
	handler := cors.Default().Handler(mux)

	// Start the server with CORS enabled
	log.Fatal(app.run(handler))
}
