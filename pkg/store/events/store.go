package events

import (
	"github.com/jinzhu/gorm"
	"github.com/rsys-speerzad/stackgen/pkg/models"
)

type Store interface {
	Create(event *models.Event) error
	Get(id string) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id string) error
	GetRecommendations(eventID string) (*models.RecommendedSlot, error)
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &store{db: db}
}

// Create inserts a new event into the database.
func (s *store) Create(event *models.Event) error {
	if err := s.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

// Get retrieves an event by its ID from the database.
func (s *store) Get(id string) (*models.Event, error) {
	var event models.Event
	if err := s.db.Preload("EventSlots,Organizer").Where("id = ?", id).First(&event).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &event, nil
}

// Update modifies an existing event in the database.
func (s *store) Update(event *models.Event) error {
	if err := s.db.Save(event).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes an event by its ID from the database.
func (s *store) Delete(id string) error {
	if err := s.db.Where("id = ?", id).Delete(&models.Event{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}
