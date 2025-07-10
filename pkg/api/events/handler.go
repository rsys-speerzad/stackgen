package events

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/rsys-speerzad/stackgen/pkg/api"
	"github.com/rsys-speerzad/stackgen/pkg/models"
	"github.com/rsys-speerzad/stackgen/pkg/store/events"
)

type Handler struct {
	store events.Store
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		store: events.NewStore(db),
	}
}

// CreateEvent handles the creation of a new event.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	var event *models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.store.Create(event); err != nil {
		http.Error(w, "Failed to create event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, event, http.StatusCreated) // Use the utility function to write the response
}

// GetEvent retrieves an event by its ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	id := urlParams.ByName("id")
	if id == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}
	event, err := h.store.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, event, 0) // Use the utility function to write the response
}

// UpdateEvent updates an existing event by its ID.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	var event *models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid resquest payload", http.StatusBadRequest)
		return
	}
	if err := h.store.Update(event); err != nil {
		http.Error(w, "Failed to create event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, event, 0) // Use the utility function to write the response
}

// DeleteEvent deletes an event by its ID.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	id := urlParams.ByName("id")
	if id == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}
	if err := h.store.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, "", 0) // Use the utility function to write the response
}

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	id := urlParams.ByName("id")
	if id == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}
	recommendations, err := h.store.GetRecommendations(id)
	if err != nil {
		http.Error(w, "Failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.ResponseWriter(w, recommendations, 0) // Use the utility function to write the response
}
