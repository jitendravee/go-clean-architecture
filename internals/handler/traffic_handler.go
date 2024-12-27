package handler

import (
	"net/http"

	"github.com/jitendravee/clean_go/internals/models"
	"github.com/jitendravee/clean_go/internals/usecase"
	"github.com/jitendravee/clean_go/internals/utils"
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
	err := utils.ReadJSON(w, r, &trafficData)
	if err != nil {
		return
	}

	createdTrafficData, err := h.trafficUserCase.Create(r.Context(), &trafficData)
	if err != nil {
		http.Error(w, "could not create the data ", http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSON(w, http.StatusCreated, createdTrafficData)
	if err != nil {
		http.Error(w, "could not write response", http.StatusInternalServerError)
	}
}
