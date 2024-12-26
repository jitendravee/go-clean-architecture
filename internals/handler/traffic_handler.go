package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jitendravee/clean_go/internals/models"
	"github.com/jitendravee/clean_go/internals/usecase"
)

type TrafficHandler struct {
	trafficUserCase *usecase.TrafficUseCase
}

func NewTrafficHandler(trafficUserCase *usecase.TrafficUseCase) *TrafficHandler {
	return &TrafficHandler{
		trafficUserCase: trafficUserCase,
	}
}

func (h *TrafficHandler) Create(w http.ResponseWriter, r *http.Request) {
	var trafficData models.Traffic
	if err := json.NewDecoder(r.Body).Decode(&trafficData); err != nil {
		http.Error(w, "could not create the jsonData ", http.StatusBadRequest)
		return
	}

	createdTrafficData, err := h.trafficUserCase.Create(r.Context(), &trafficData)
	if err != nil {
		http.Error(w, "could not create the data ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTrafficData)
}
