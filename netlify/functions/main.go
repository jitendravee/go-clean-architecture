// // // package main

// // // import (
// // // 	"context"
// // // 	"fmt"
// // // 	"log"
// // // 	"net/http"
// // // 	"os"
// // // 	"time"

// // // 	"github.com/aws/aws-lambda-go/events"
// // // 	"github.com/aws/aws-lambda-go/lambda"
// // // 	"github.com/go-chi/chi/v5"
// // // 	"github.com/go-chi/chi/v5/middleware"
// // // 	"github.com/jitendravee/clean_go/internals/db"
// // // 	"github.com/jitendravee/clean_go/internals/repository"
// // // 	"github.com/jitendravee/clean_go/internals/store"
// // // 	"github.com/jitendravee/clean_go/internals/usecase"
// // // 	"github.com/joho/godotenv"
// // // 	"go.mongodb.org/mongo-driver/mongo"
// // // )

// // // type dbConfig struct {
// // // 	addr   string
// // // 	dbName string
// // // }

// // // type config struct {
// // // 	db dbConfig
// // // }

// // // type application struct {
// // // 	config config
// // // 	store  store.Storage
// // // }

// // // func loadConfig() (config, error) {
// // // 	err := godotenv.Load()
// // // 	if err != nil {
// // // 		log.Println("Error loading .env file, proceeding with default environment variables")
// // // 	}

// // // 	mongoURL := os.Getenv("MONGO_URL")
// // // 	if mongoURL == "" {
// // // 		return config{}, fmt.Errorf("MONGO_URL environment variable not set")
// // // 	}

// // // 	dbName := os.Getenv("DB_NAME")
// // // 	if dbName == "" {
// // // 		dbName = "user"
// // // 	}

// // // 	return config{
// // // 		db: dbConfig{
// // // 			addr:   mongoURL,
// // // 			dbName: dbName,
// // // 		},
// // // 	}, nil
// // // }
// // // func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// // // 	cfg, err := loadConfig()
// // // 	if err != nil {
// // // 		return events.APIGatewayProxyResponse{
// // // 			StatusCode: 500,
// // // 			Body:       fmt.Sprintf("Error loading configuration: %v", err),
// // // 		}, nil
// // // 	}

// // // 	database, err := db.New(cfg.db.addr, cfg.db.dbName)
// // // 	if err != nil {
// // // 		return events.APIGatewayProxyResponse{
// // // 			StatusCode: 500,
// // // 			Body:       fmt.Sprintf("Error connecting to the database: %v", err),
// // // 		}, nil
// // // 	}

// // // 	app := &application{
// // // 		config: cfg,
// // // 		store:  store.NewStorage(database),
// // // 	}

// // // func (app *application) mount(db *mongo.Database) http.Handler {
// // // 	trafficRepo := repository.NewMongoTrafficRepo(db)
// // // 	trafficUseCase := usecase.NewTrafficUseCase(trafficRepo)
// // // 	trafficHandler := handler.NewTrafficHandler(trafficUseCase)

// // // 	r := chi.NewRouter()
// // // 	r.Use(middleware.RequestID)
// // // 	r.Use(middleware.RealIP)
// // // 	r.Use(middleware.Logger)
// // // 	r.Use(middleware.Recoverer)
// // // 	r.Use(middleware.Timeout(60 * time.Second))

// // // 	r.Route("/v1", func(r chi.Router) {
// // // 		r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
// // // 			w.Write([]byte("hello"))
// // // 		})
// // // 		r.Post("/traffic", trafficHandler.Create)
// // // 	})

// // // 	return r
// // // }

// // // 	router := app.mount(database)
// // // 	reqCtx := context.Background()
// // // 	w := newLambdaResponseWriter()
// // // 	r, err := events.APIGatewayProxyRequestToHttpRequest(reqCtx, req)
// // // 	if err != nil {
// // // 		return events.APIGatewayProxyResponse{
// // // 			StatusCode: 500,
// // // 			Body:       fmt.Sprintf("Error converting request: %v", err),
// // // 		}, nil
// // // 	}

// // // 	router.ServeHTTP(w, r)

// // // 	return w.ToAPIGatewayProxyResponse(), nil
// // // }

// // //	func main() {
// // //		lambda.Start(handler)
// // //	}
// // package main

// // import (
// // 	"context"
// // 	"fmt"
// // 	"log"
// // 	"net/http"
// // 	"os"
// // 	"time"

// // 	"github.com/aws/aws-lambda-go/events"
// // 	"github.com/aws/aws-lambda-go/lambda"
// // 	"github.com/go-chi/chi/v5"
// // 	"github.com/go-chi/chi/v5/middleware"
// // 	"github.com/jitendravee/clean_go/internals/db"
// // 	"github.com/jitendravee/clean_go/internals/handler"
// // 	"github.com/jitendravee/clean_go/internals/repository"
// // 	"github.com/jitendravee/clean_go/internals/store"
// // 	"github.com/jitendravee/clean_go/internals/usecase"
// // 	"github.com/joho/godotenv"
// // 	"go.mongodb.org/mongo-driver/mongo"
// // )

// // type dbConfig struct {
// // 	addr   string
// // 	dbName string
// // }

// // type config struct {
// // 	db dbConfig
// // }

// // type application struct {
// // 	config config
// // 	store  store.Storage
// // }

// // func loadConfig() (config, error) {
// // 	err := godotenv.Load()
// // 	if err != nil {
// // 		log.Println("Error loading .env file, proceeding with default environment variables")
// // 	}

// // 	mongoURL := os.Getenv("MONGO_URL")
// // 	if mongoURL == "" {
// // 		return config{}, fmt.Errorf("MONGO_URL environment variable not set")
// // 	}

// // 	dbName := os.Getenv("DB_NAME")
// // 	if dbName == "" {
// // 		dbName = "user"
// // 	}

// // 	return config{
// // 		db: dbConfig{
// // 			addr:   mongoURL,
// // 			dbName: dbName,
// // 		},
// // 	}, nil
// // }

// // func (app *application) mount(db *mongo.Database) http.Handler {
// // 	trafficRepo := repository.NewMongoTrafficRepo(db)
// // 	trafficUseCase := usecase.NewTrafficUseCase(trafficRepo)
// // 	trafficHandler := handler.NewTrafficHandler(trafficUseCase)

// // 	r := chi.NewRouter()
// // 	r.Use(middleware.RequestID)
// // 	r.Use(middleware.RealIP)
// // 	r.Use(middleware.Logger)
// // 	r.Use(middleware.Recoverer)
// // 	r.Use(middleware.Timeout(60 * time.Second))

// // 	r.Route("/v1", func(r chi.Router) {
// // 		r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
// // 			w.Write([]byte("hello"))
// // 		})
// // 		r.Post("/traffic", trafficHandler.Create)
// // 	})

// // 	return r
// // }

// // func ahandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// // 	cfg, err := loadConfig()
// // 	if err != nil {
// // 		return events.APIGatewayProxyResponse{
// // 			StatusCode: 500,
// // 			Body:       fmt.Sprintf("Error loading configuration: %v", err),
// // 		}, nil
// // 	}

// // 	database, err := db.New(cfg.db.addr, cfg.db.dbName)
// // 	if err != nil {
// // 		return events.APIGatewayProxyResponse{
// // 			StatusCode: 500,
// // 			Body:       fmt.Sprintf("Error connecting to the database: %v", err),
// // 		}, nil
// // 	}

// // 	app := &application{
// // 		config: cfg,
// // 		store:  store.NewStorage(database),
// // 	}

// // 	router := app.mount(database)
// // 	reqCtx := context.Background()
// // 	w := newLambdaResponseWriter()
// // 	r, err := events.APIGatewayProxyRequestToHttpRequest(reqCtx, req)
// // 	if err != nil {
// // 		return events.APIGatewayProxyResponse{
// // 			StatusCode: 500,
// // 			Body:       fmt.Sprintf("Error converting request: %v", err),
// // 		}, nil
// // 	}

// // 	router.ServeHTTP(w, r)

// // 	return w.ToAPIGatewayProxyResponse(), nil
// // }

// // func main() {
// // 	lambda.Start(ahandler)
// // }
// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/jitendravee/clean_go/internals/db"
// 	"github.com/jitendravee/clean_go/internals/handler"
// 	"github.com/jitendravee/clean_go/internals/repository"
// 	"github.com/jitendravee/clean_go/internals/store"
// 	"github.com/jitendravee/clean_go/internals/usecase"
// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type dbConfig struct {
// 	addr   string
// 	dbName string
// }

// type config struct {
// 	db dbConfig
// }

// type application struct {
// 	config config
// 	store  store.Storage
// }

// func loadConfig() (config, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("Error loading .env file, proceeding with default environment variables")
// 	}

// 	mongoURL := os.Getenv("MONGO_URL")
// 	if mongoURL == "" {
// 		return config{}, fmt.Errorf("MONGO_URL environment variable not set")
// 	}

// 	dbName := os.Getenv("DB_NAME")
// 	if dbName == "" {
// 		dbName = "user"
// 	}

// 	return config{
// 		db: dbConfig{
// 			addr:   mongoURL,
// 			dbName: dbName,
// 		},
// 	}, nil
// }

// // newLambdaResponseWriter creates a custom response writer for AWS Lambda
// type lambdaResponseWriter struct {
// 	statusCode int
// 	body       string
// 	headers    map[string][]string
// }

// func newLambdaResponseWriter() *lambdaResponseWriter {
// 	return &lambdaResponseWriter{
// 		headers: make(map[string][]string),
// 	}
// }

// func (w *lambdaResponseWriter) Header() http.Header {
// 	return w.headers
// }

// func (w *lambdaResponseWriter) Write(b []byte) (int, error) {
// 	w.body = string(b)
// 	return len(b), nil
// }

// func (w *lambdaResponseWriter) WriteHeader(statusCode int) {
// 	w.statusCode = statusCode
// }

// func (w *lambdaResponseWriter) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: w.statusCode,
// 		Body:       w.body,
// 		Headers:    w.headers,
// 	}
// }

// func (app *application) mount(db *mongo.Database) http.Handler {
// 	trafficRepo := repository.NewMongoTrafficRepo(db)
// 	trafficUseCase := usecase.NewTrafficUseCase(trafficRepo)
// 	trafficHandler := handler.NewTrafficHandler(trafficUseCase)

// 	r := chi.NewRouter()
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)
// 	r.Use(middleware.Timeout(60 * time.Second))

// 	r.Route("/v1", func(r chi.Router) {
// 		r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte("hello"))
// 		})
// 		r.Post("/traffic", trafficHandler.Create)
// 	})

// 	return r
// }

// func ahandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	cfg, err := loadConfig()
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 500,
// 			Body:       fmt.Sprintf("Error loading configuration: %v", err),
// 		}, nil
// 	}

// 	database, err := db.New(cfg.db.addr, cfg.db.dbName)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 500,
// 			Body:       fmt.Sprintf("Error connecting to the database: %v", err),
// 		}, nil
// 	}

// 	app := &application{
// 		config: cfg,
// 		store:  store.NewStorage(database),
// 	}

// 	router := app.mount(database)

// 	// Convert APIGatewayProxyRequest to http.Request manually
// 	headers := make(map[string][]string)
// 	for key, value := range req.Headers {
// 		headers[key] = append(headers[key], value)
// 	}

// 	// Create the HTTP request
// 	r, err := http.NewRequest(req.HTTPMethod, req.Path, strings.NewReader(req.Body))
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 500,
// 			Body:       fmt.Sprintf("Error converting request: %v", err),
// 		}, nil
// 	}
// 	r.Header = headers

// 	w := newLambdaResponseWriter()

// 	// Call the router's ServeHTTP method to process the request
// 	router.ServeHTTP(w, r)

// 	return w.ToAPIGatewayProxyResponse(), nil
// }

//	func main() {
//		lambda.Start(ahandler)
//	}
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	db dbConfig
}

type application struct {
	config config
	store  store.Storage
}

func loadConfig() (config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, proceeding with default environment variables")
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return config{}, fmt.Errorf("MONGO_URL environment variable not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "user"
	}

	return config{
		db: dbConfig{
			addr:   mongoURL,
			dbName: dbName,
		},
	}, nil
}

// newLambdaResponseWriter creates a custom response writer for AWS Lambda
type lambdaResponseWriter struct {
	statusCode int
	body       string
	headers    map[string]string // Adjusted to be map[string]string
}

func newLambdaResponseWriter() *lambdaResponseWriter {
	return &lambdaResponseWriter{
		headers: make(map[string]string),
	}
}

func (w *lambdaResponseWriter) Header() http.Header {
	// Convert map[string]string to http.Header
	headers := http.Header{}
	for key, value := range w.headers {
		headers.Add(key, value)
	}
	return headers
}

func (w *lambdaResponseWriter) Write(b []byte) (int, error) {
	w.body = string(b)
	return len(b), nil
}

func (w *lambdaResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *lambdaResponseWriter) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	// Convert headers to map[string]string for APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		StatusCode: w.statusCode,
		Body:       w.body,
		Headers:    w.headers,
	}
}

func (app *application) mount(db *mongo.Database) http.Handler {
	trafficRepo := repository.NewMongoTrafficRepo(db)
	trafficUseCase := usecase.NewTrafficUseCase(trafficRepo)
	trafficHandler := handler.NewTrafficHandler(trafficUseCase)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})
		r.Post("/traffic", trafficHandler.Create)
	})

	return r
}

func ahandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := loadConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error loading configuration: %v", err),
		}, nil
	}

	database, err := db.New(cfg.db.addr, cfg.db.dbName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error connecting to the database: %v", err),
		}, nil
	}

	app := &application{
		config: cfg,
		store:  store.NewStorage(database),
	}

	router := app.mount(database)

	// Convert APIGatewayProxyRequest to http.Request manually
	headers := make(map[string][]string)
	for key, value := range req.Headers {
		headers[key] = append(headers[key], value)
	}

	// Create the HTTP request
	r, err := http.NewRequest(req.HTTPMethod, req.Path, strings.NewReader(req.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error converting request: %v", err),
		}, nil
	}
	r.Header = headers

	w := newLambdaResponseWriter()

	// Call the router's ServeHTTP method to process the request
	router.ServeHTTP(w, r)

	return w.ToAPIGatewayProxyResponse(), nil
}

func main() {
	lambda.Start(ahandler)
}
