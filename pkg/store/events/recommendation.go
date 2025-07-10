package events

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
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
	if err := s.db.Preload("Availabilities", "event_id IS NULL and start_time >= ?", event.EventSlots[0].StartTime).Find(&users).Error; err != nil {
		fmt.Println("Error retrieving users:", err)
		return nil, err
	}
	slots := make(map[uuid.UUID][][]string)
	for _, user := range users {
		// check if user has any availability for given the event slot
		if len(user.Availabilities) == 0 {
			continue
		}
		consecutiveAvailabilities := MergeConsecutiveAvailabilities(user.Availabilities)
		for _, slot := range event.EventSlots {
			available := false
			for _, availability := range consecutiveAvailabilities {
				if slot.StartTime.Before(availability.StartTime) && slot.EndTime.After(availability.EndTime) ||
					slot.StartTime.Equal(availability.StartTime) && slot.EndTime.Equal(availability.EndTime) {
					available = true
					slots[slot.ID][0] = append(slots[slot.ID][0], user.ID.String())
					break
				}
			}
			if !available {
				slots[slot.ID] = append(slots[slot.ID], []string{user.ID.String()})
			}
		}
	}

	// for each event slot, check which users are available, and which are not
	// finalSlot := event.EventSlots[0]
	// availableUser, missingUsers := checkUserAvailability(userAvailabilities, &finalSlot)
	// for i := 1; i < len(event.EventSlots); i++ {
	// 	currSlotAvailableUser, currSlotMissingUsers := checkUserAvailability(userAvailabilities, &event.EventSlots[i])
	// 	if len(currSlotAvailableUser) > len(availableUser) {
	// 		availableUser = currSlotAvailableUser
	// 		missingUsers = currSlotMissingUsers
	// 		finalSlot = event.EventSlots[i]
	// 	}
	// }
	// return the slot with the most available users
	return &models.RecommendedSlot{
		// StartTime:      finalSlot.StartTime,
		// EndTime:        finalSlot.EndTime,
		// UserIDs:        availableUser,
		// MissingUserIDs: missingUsers,
	}, nil
}

func checkUserAvailability(availabilities []*models.UserAvailability, slot *models.EventSlot) ([]string, []string) {
	availableUsers := make(map[string]struct{})
	missingUsers := make(map[string]struct{})
	for _, availability := range availabilities {
		if slot.StartTime.Before(availability.StartTime) && slot.EndTime.After(availability.EndTime) ||
			slot.StartTime.Equal(availability.StartTime) && slot.EndTime.Equal(availability.EndTime) {
			availableUsers[availability.UserID.String()] = struct{}{}
			continue
		}
		missingUsers[availability.UserID.String()] = struct{}{}

	}
	var availableUserIDs []string
	var missingUserIDs []string
	for userID := range availableUsers {
		availableUserIDs = append(availableUserIDs, userID)
	}
	for userID := range missingUsers {
		missingUserIDs = append(missingUserIDs, userID)
	}
	return availableUserIDs, missingUserIDs
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
