package events

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func InitializeRouter(r *httprouter.Router, db *gorm.DB) {
	handler := NewHandler(db)
	r.POST("/event", handler.Create)                                 // Create a new event
	r.GET("/event/:id", handler.Get)                                 // Get event by ID
	r.PUT("/event/:id", handler.Update)                              // Update event by ID
	r.DELETE("/event/:id", handler.Delete)                           // Delete event by ID
	r.GET("/events/:id/recommendations", handler.GetRecommendations) // Get recommendations for an event
}
