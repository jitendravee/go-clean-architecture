package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jitendravee/clean_go/internals/handler"
	"github.com/jitendravee/clean_go/internals/repository"
	"github.com/jitendravee/clean_go/internals/store"
	"github.com/jitendravee/clean_go/internals/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

type config struct {
	addr string
	db   dbConfig
}

type application struct {
	config config
	store  store.Storage
}
type dbConfig struct {
	addr   string
	dbName string
}

func (app *application) mount(db *mongo.Database) http.Handler {
	trafficRepo := repository.NewMongoTrafficRepo(db)
	trafficUserCase := usecase.NewTrafficUseCase(trafficRepo)
	trafficHandler := handler.NewTrafficHandler(trafficUserCase)
	signalRepo := repository.NewSignalRepo(db)
	signalUseCase := usecase.NewSignalUseCase(signalRepo)
	signalHandler := handler.NewSignalHandler(signalUseCase)

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
		r.Post("/user", app.InsertUser)
		r.Post("/traffic", trafficHandler.Create)
		r.Post("/signal", signalHandler.CreateSignal)
		r.Get("/signal", signalHandler.GetSignalHandler)
		r.Get("/signal/{group_id}", signalHandler.GetGroupSignalByIdHandler)
		r.Patch("/signal/{group_id}/update-count", signalHandler.UpdateVechileCountHandler)
		r.Patch("/signal/{group_id}/update-image", signalHandler.UpdateImageUrlHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	log.Printf("the server is running on port : %s", app.config.addr)
	return srv.ListenAndServe()

}
