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
	// Step 1: Get the ID from the URL parameters
	groupId := utils.GetParams(r, "group_id")

	// Step 2: Check if the ID is valid
	if groupId == "" {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	// Step 3: Call the use case to get the group signal data
	groupSignalData, err := h.signalUseCase.GetGroupSignalByIdUseCase(r.Context(), groupId)
	if err != nil {
		// Return a more descriptive error if the use case failed
		http.Error(w, fmt.Sprintf("Error retrieving group signal data: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 4: Check if the data is nil or empty
	if groupSignalData == nil {
		// If no data is found, return 404 (Not Found)
		http.Error(w, "Group signal not found", http.StatusNotFound)
		return
	}

	// Step 5: Write the JSON response
	err = utils.WriteJSON(w, http.StatusOK, groupSignalData)
	if err != nil {
		// Handle error if JSON writing fails
		fmt.Printf("Error writing the JSON response: %v\n", err)
		http.Error(w, "Could not write JSON response", http.StatusInternalServerError)
	}
}
