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
		http.Error(w, fmt.Sprintf("could not get signals: %v", err), http.StatusInternalServerError)

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
func (h *SignalHandler) GetSignalHandler(w http.ResponseWriter, r *http.Request) {
	signalsData, err := h.signalUseCase.SignalRepo.GetAllSignal(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get signals: %v", err), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSON(w, http.StatusOK, signalsData)
	if err != nil {
		fmt.Printf("Error writing JSON response: %v\n", err)
		http.Error(w, "could not write response", http.StatusInternalServerError)
	}
}

func (h *SignalHandler) GetGroupSignalByIdHandler(w http.ResponseWriter, r *http.Request) {
	groupId := utils.GetParams(r, "group_id")

	if groupId == "" {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	groupSignalData, err := h.signalUseCase.GetGroupSignalByIdUseCase(r.Context(), groupId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving group signal data: %v", err), http.StatusInternalServerError)
		return
	}

	if groupSignalData == nil {
		http.Error(w, "Group signal not found", http.StatusNotFound)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, groupSignalData)
	if err != nil {
		fmt.Printf("Error writing the JSON response: %v\n", err)
		http.Error(w, "Could not write JSON response", http.StatusInternalServerError)
	}
}

func (h *SignalHandler) UpdateVechileCountHandler(w http.ResponseWriter, r *http.Request) {
	groupId := utils.GetParams(r, "group_id")
	var updateCountRequest models.UpdateSignalCountGroup
	err := utils.ReadJSON(w, r, &updateCountRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	updateSignalData, err := h.signalUseCase.UpdateVechileCountBySignalIdUseCase(r.Context(), &updateCountRequest, groupId)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update vehicle count: %v", err), http.StatusInternalServerError)
		return

	}
	err = utils.WriteJSON(w, http.StatusOK, updateSignalData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)

	}
}
func (h *SignalHandler) UpdateImageUrlHandler(w http.ResponseWriter, r *http.Request) {
	groupId := utils.GetParams(r, "group_id")
	var updateImageUrlRequest models.ImageRequestList
	err := utils.ReadJSON(w, r, &updateImageUrlRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	updatedImageData, err := h.signalUseCase.UpdateImageUrlUseCase(r.Context(), &updateImageUrlRequest, groupId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update vehicle count: %v", err), http.StatusInternalServerError)
		return

	}
	err = utils.WriteJSON(w, http.StatusOK, updatedImageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)

	}

}
