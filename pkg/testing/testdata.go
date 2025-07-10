package testing

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/rsys-speerzad/stackgen/pkg/models"
)

// CreateTestUsersData creates users with availabilities and random events
func CreateTestData(db *gorm.DB) error {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	for i := 1; i <= 10; i++ {
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC) // Set a fixed start time for all users
		user := models.User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d_%d@example.com", i, rand.Intn(10000)),
		}
		for j := 0; j < 20; j++ { // Each user has 20 availabilities from 9 AM to 6 PM of 30 mins each
			user.Availabilities = append(user.Availabilities, &models.UserAvailability{
				Slot: models.Slot{
					StartTime: startTime,
					EndTime:   startTime.Add(30 * time.Minute),
				},
			})
			startTime = startTime.Add(30 * time.Minute) // Increment start time for next availability
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}
	eventIDs := []uuid.UUID{}
	for i := 1; i <= 5; i++ { // Create 5 random events
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local)
		event := models.Event{
			Title:             fmt.Sprintf("Event%d", i),
			Description:       fmt.Sprintf("Description for Event%d", i),
			EstimatedDuration: rand.Intn(120) + 30, // Random duration between 30 and 150 minutes
			OrganizerID:       nil,                 // No organizer for now
		}
		for j := 0; j < 3; j++ { // Each event has
			event.EventSlots = append(event.EventSlots, models.EventSlot{
				StartTime: startTime,
				EndTime: startTime.Add(time.Duration(rand.Intn(120)+30) *
					time.Minute), // Random slot duration between 30 and 150 minutes
			})
			startTime = startTime.Add(time.Duration(rand.Intn(60)+30) * time.Minute) // Increment start time for next slot
		}
		if err := db.Create(&event).Error; err != nil {
			return err
		}
		eventIDs = append(eventIDs, event.ID)
	}
	var users []models.User
	if err := db.Preload("Availabilities").Order("RANDOM()").Limit(3).Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		for _, avail := range user.Availabilities {
			if len(eventIDs) == 0 {
				break
			}
			// Assign a random event ID to the availability
			randomEventID := eventIDs[rand.Intn(len(eventIDs))]
			avail.EventID = &randomEventID
			if err := db.Model(avail).Update("event_id", randomEventID).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
