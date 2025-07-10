package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/rsys-speerzad/stackgen/pkg/api"
	"github.com/rsys-speerzad/stackgen/pkg/models"
	"github.com/rsys-speerzad/stackgen/pkg/store/users"
)

type Handler struct {
	store users.Store
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		store: users.NewStore(db),
	}
}

// Create creates a new user
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.store.Create(user); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, user, http.StatusCreated) // Use the utility function to write the response
}

// Get gets the user by ID
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	id := urlParams.ByName("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	user, err := h.store.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, user, 0) // Use the utility function to write the response
}

// Update updates the user by ID
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	id := urlParams.ByName("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.store.Update(id, user); err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, "user updated successfully", 0) // Use the utility function to write the response
}

// Delete deletes the user by ID
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	id := urlParams.ByName("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if err := h.store.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, "user deleted successfully", http.StatusNoContent) // No content to return
}

func (h *Handler) AddAvailability(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// read user ID from URL parameters
	UserID := urlParams.ByName("id")
	if UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	// parse the user ID to uuid.UUID
	userID, err := uuid.Parse(UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	// decode the request body to get availability slots
	var req = struct {
		Slots []models.Slot `json:"slots"`
	}{}
	// check if the request body is valid
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("error while decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// convert the slots to UserAvailability models
	var slots []models.UserAvailability
	for _, slot := range req.Slots {
		slots = append(slots, models.UserAvailability{
			UserID: userID,
			Slot: models.Slot{
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
			},
		})
	}
	// add the availability slots to the user
	if err := h.store.AddAvailability(slots); err != nil {
		http.Error(w, "Failed to add availability: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, nil, http.StatusNoContent) // No content to return
}

func (h *Handler) UpdateAvailability(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// read user ID from URL parameters
	UserID := urlParams.ByName("id")
	if UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	// parse the user ID to uuid.UUID
	userID, err := uuid.Parse(UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	// read availability ID from URL parameters
	id := urlParams.ByName("aid")
	if id == "" {
		http.Error(w, "Availability ID is required", http.StatusBadRequest)
		return
	}
	// decode the request body to get updated availability slot
	var slot models.UserAvailability
	if err := json.NewDecoder(r.Body).Decode(&slot); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	slot.UserID = userID // set the user ID
	if err := h.store.UpdateAvailability(id, slot); err != nil {
		http.Error(w, "Failed to update availability: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, nil, http.StatusNoContent) // No content to return
}

func (h *Handler) DeleteAvailability(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// read availability ID from URL parameters
	id := urlParams.ByName("aid")
	if id == "" {
		http.Error(w, "Availability ID is required", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteAvailability(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Availability not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete availability: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, nil, http.StatusNoContent) // No content to return
}

func (h *Handler) GetAvailabilities(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	// read user ID from URL parameters
	userID := urlParams.ByName("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	// get the availabilities for the user
	availabilities, err := h.store.GetAvailability(userID)
	if err != nil {
		http.Error(w, "Failed to get availabilities: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var resp = struct {
		Slots []models.Slot `json:"available_slots"`
	}{}
	for _, availability := range availabilities {
		resp.Slots = append(resp.Slots, models.Slot{
			StartTime: availability.StartTime,
			EndTime:   availability.EndTime,
		})
	}
	api.ResponseWriter(w, resp, 0) // Use the utility function to write the response
}
