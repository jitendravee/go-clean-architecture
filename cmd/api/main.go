package main

import (
	"log"

	"github.com/jitendravee/clean_go/internals/db"
	"github.com/jitendravee/clean_go/internals/store"
)

func main() {
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			addr:   "mongodb+srv://jitendrajat6397:Jitendra6323@jitendra.yaofk.mongodb.net/",
			dbName: "user",
		},
	}
	db, err := db.New(
		cfg.db.addr,
		cfg.db.dbName,
	)
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
