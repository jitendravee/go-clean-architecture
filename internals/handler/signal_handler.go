package handler

import (
	"fmt"
	"net/http"

	"github.com/jitendravee/clean_go/internals/models"
	"github.com/jitendravee/clean_go/internals/usecase"
	"github.com/jitendravee/clean_go/internals/utils"
)

type SignalHandler struct {
	signalUseCase *usecase.SignalUseCase
}

func NewSignalHandler(signalUseCase *usecase.SignalUseCase) *SignalHandler {
	return &SignalHandler{
		signalUseCase: signalUseCase,
	}

}

func (h *SignalHandler) CreateSignal(w http.ResponseWriter, r *http.Request) {
	var signalData models.GroupSignal
	err := utils.ReadJSON(w, r, &signalData)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		fmt.Printf("Error reading JSON: %v\n", err)
		return
	}
	fmt.Printf("Decoded signal data: %+v\n", signalData)

	createdSignalData, err := h.signalUseCase.CreateGroupSignal(r.Context(), &signalData)
	if err != nil {
		fmt.Printf("Error creating signal: %v\n", err)
		http.Error(w, "could not create the data", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Created signal data: %+v\n", createdSignalData)

	err = utils.WriteJSON(w, http.StatusCreated, createdSignalData)
	if err != nil {
		fmt.Printf("Error writing JSON response: %v\n", err)
		http.Error(w, "could not write response", http.StatusInternalServerError)
	}
}
