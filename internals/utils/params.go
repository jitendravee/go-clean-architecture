package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetParams(r *http.Request, param string) string {
	value := chi.URLParam(r, param)
	return value
}
