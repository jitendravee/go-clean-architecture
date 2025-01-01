package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jitendravee/clean_go/internals/db"
	"github.com/jitendravee/clean_go/internals/handler"
	"github.com/jitendravee/clean_go/internals/repository"
	"github.com/jitendravee/clean_go/internals/store"
	"github.com/jitendravee/clean_go/internals/usecase"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type dbConfig struct {
	addr   string
	dbName string
}

type config struct {
	addr string
	db   dbConfig
}

type application struct {
	config config
	store  store.Storage
}

// Load environment variables from .env file and initialize the config
func loadConfig() (config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080" // Default address if not set
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return config{}, fmt.Errorf("MONGO_URL environment variable not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "user" // Default database name
	}

	return config{
		addr: addr,
		db: dbConfig{
			addr:   mongoURL,
			dbName: dbName,
		},
	}, nil
}

// Initialize the application with routing and handlers
func (app *application) mount(db *mongo.Database) http.Handler {
	// Create repositories, use cases, and handlers
	trafficRepo := repository.NewMongoTrafficRepo(db)
	trafficUseCase := usecase.NewTrafficUseCase(trafficRepo)
	trafficHandler := handler.NewTrafficHandler(trafficUseCase)

	signalRepo := repository.NewSignalRepo(db)
	signalUseCase := usecase.NewSignalUseCase(signalRepo)
	signalHandler := handler.NewSignalHandler(signalUseCase)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Define routes
	r.Route("/v1", func(r chi.Router) {
		r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})
		r.Post("/traffic", trafficHandler.Create)
		r.Post("/signal", signalHandler.CreateSignal)
		r.Get("/signal", signalHandler.GetSignalHandler)
		r.Get("/signal/{group_id}", signalHandler.GetGroupSignalByIdHandler)
		r.Patch("/signal/{group_id}/update-count", signalHandler.UpdateVechileCountHandler)
	})

	return r
}

// Start the application server
func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is running on port %s", app.config.addr)
	return srv.ListenAndServe()
}

// Main handler for Netlify function
func Handler(w http.ResponseWriter, r *http.Request) {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	// Initialize the database connection
	db, err := db.New(cfg.db.addr, cfg.db.dbName)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Initialize the storage and application
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	// Mount routes
	mux := app.mount(db)

	// Start the application server
	if err := app.run(mux); err != nil {
		log.Fatal("Error running server: ", err)
	}

	// Send a success response
	fmt.Fprintf(w, "App is running successfully!")
}

// Main entry point for the application
func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
