// package main

// import (
// 	"flag"
// 	"log"
// 	"os"

// 	"github.com/jitendravee/clean_go/internals/db"
// 	"github.com/jitendravee/clean_go/internals/store"
// 	"github.com/joho/godotenv"
// 	"github.com/rs/cors"
// )

// // func main() {

// // 	// envPath := flag.String("env", ".env", "Path to .env file")
// // 	flag.Parse()

// // 	// Load .env file
// // 	err := godotenv.Load("/home/jitendra/Downloads/go_project/clean_go/cmd/api/.env")
// // 	if err != nil {
// // 		log.Fatal("Error loading .env file")
// // 	}

// // 	// Get environment variables
// // 	addr := os.Getenv("ADDR")
// // 	if addr == "" {
// // 		addr = ":8080" // Default port
// // 	}

// // 	mongoURL := os.Getenv("MONGO_URL")
// // 	if mongoURL == "" {
// // 		log.Fatal("MONGO_URL environment variable not set")
// // 	}

// // 	dbName := os.Getenv("DB_NAME")
// // 	if dbName == "" {
// // 		dbName = "user" // Default database name
// // 	}

// // 	// Set up configuration
// // 	cfg := config{
// // 		addr: addr,
// // 		db: dbConfig{
// // 			addr:   mongoURL,
// // 			dbName: dbName,
// // 		},
// // 	}

// // 	// Initialize database connection
// // 	db, err := db.New(cfg.db.addr, cfg.db.dbName)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	store := store.NewStorage(db)

// // 	// Set up application and routes
// // 	app := &application{
// // 		config: cfg,
// // 		store:  store,
// // 	}
// // 	mux := app.mount(db)

// // 	// Enable CORS with middleware
// // 	handler := cors.Default().Handler(mux) // Default CORS settings allow all origins

// // 	// Start the server with CORS enabled
// // 	log.Fatal(app.run(handler))
// // }

package main

import (
	"log"
	"os"

	"github.com/jitendravee/clean_go/internals/db"
	"github.com/jitendravee/clean_go/internals/store"
	"github.com/rs/cors"
)

func main() {
	// Get environment variables
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

	// Customize CORS configuration to allow all methods, headers, and credentials
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true, // Allow credentials (cookies, authorization headers, etc.)
	})

	// Wrap the mux with the CORS middleware
	handler := c.Handler(mux)

	// Start the server
	log.Fatal(app.run(handler))
}
