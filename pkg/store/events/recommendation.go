package events

import (
	"fmt"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/rsys-speerzad/stackgen/pkg/models"
)

// GetRecommendations retrieves recommended slots for an event based on user availability.
func (s *store) GetRecommendations(eventID string) (*models.RecommendedSlot, error) {
	var event models.Event
	var users []models.User
	if err := s.db.Preload("EventSlots").First(&event, "id = ?", eventID).Error; err != nil {
		fmt.Println("Error retrieving event:", err)
		if gorm.IsRecordNotFoundError(err) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if err := s.db.Preload("Availabilities", "event_id IS NULL").Find(&users).Error; err != nil {
		fmt.Println("Error retrieving users:", err)
		return nil, err
	}
	// Logic to calculate recommended slots based on event and user availability
	// This is a placeholder; actual implementation would involve more complex logic
	var finalSlot models.EventSlot
	var finalAvailableUsers, finalMissingUsers []string
	for i := 0; i < len(event.EventSlots); i++ {
		var availableUsers, missingUsers []string
		for _, user := range users {
			mergedAvailabilities := MergeConsecutiveAvailabilities(user.Availabilities)
			if checkAvailability(event.EventSlots[i], mergedAvailabilities) {
				availableUsers = append(availableUsers, user.Name)
			} else {
				missingUsers = append(missingUsers, user.Name)
			}
		}
		if len(availableUsers) > len(finalAvailableUsers) {
			finalSlot = event.EventSlots[i]
			finalAvailableUsers = availableUsers
			finalMissingUsers = missingUsers
		}
	}
	return &models.RecommendedSlot{
		StartTime:      finalSlot.StartTime,
		EndTime:        finalSlot.EndTime,
		UserIDs:        finalAvailableUsers,
		MissingUserIDs: finalMissingUsers,
	}, nil
}

func MergeConsecutiveAvailabilities(slots []*models.UserAvailability) []*models.UserAvailability {
	if len(slots) == 0 {
		return nil
	}

	// Sort by start time
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].StartTime.Before(slots[j].StartTime)
	})

	var merged []*models.UserAvailability
	current := slots[0]

	for i := 1; i < len(slots); i++ {
		next := slots[i]
		// If next starts at or before current ends, or exactly when it ends
		if !next.StartTime.After(current.EndTime) {
			// Merge the two
			if next.EndTime.After(current.EndTime) {
				current.EndTime = next.EndTime
			}
		} else {
			merged = append(merged, current)
			current = next
		}
	}

	// Add the last one
	merged = append(merged, current)
	return merged
}

func checkAvailability(slot models.EventSlot, availabilities []*models.UserAvailability) bool {
	for _, availability := range availabilities {
		if slot.StartTime.Before(availability.EndTime) && slot.EndTime.After(availability.StartTime) {
			return true
		}
	}
	return false
}
