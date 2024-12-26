package main

import (
	"encoding/json"
	"net/http"

	"github.com/jitendravee/clean_go/internals/models"
)

func (app *application) InsertUser(w http.ResponseWriter, r *http.Request) {
	var userReq models.User

	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {

		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	ctx := r.Context()

	createdUser, err := app.store.Users.Create(ctx, user)
	if err != nil {

		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {

		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
