package users

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func InitializeRouter(r *httprouter.Router, db *gorm.DB) {
	handler := NewHandler(db)
	r.POST("/user", handler.Create)       // Create a new user
	r.GET("/user/:id", handler.Get)       // Get user by ID
	r.PUT("/user/:id", handler.Update)    // Update user by ID
	r.DELETE("/user/:id", handler.Delete) // Delete user by ID

	// availability routes
	r.GET("/user/:id/availability", handler.GetAvailabilities)          // Get availability for user
	r.POST("/user/:id/availability", handler.AddAvailability)           // Add availability for user
	r.PUT("/user/:id/availability/:aid", handler.UpdateAvailability)    // Update availability for user
	r.DELETE("/user/:id/availability/:aid", handler.DeleteAvailability) // Delete availability for user
}
