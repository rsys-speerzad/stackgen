package users

import (
	"github.com/jinzhu/gorm"
	"github.com/rsys-speerzad/stackgen/pkg/models"
)

type Store interface {
	Create(user *models.User) error
	Get(id string) (*models.User, error)
	Update(id string, user *models.User) error
	Delete(id string) error
	GetAvailability(userID string) ([]models.UserAvailability, error)
	AddAvailability(slots []models.UserAvailability) error
	UpdateAvailability(slotID string, slot models.UserAvailability) error
	DeleteAvailability(slotID string) error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &store{db: db}
}

func (s *store) Create(user *models.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) Get(id string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Availabilities").Where("id = ?", id).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *store) Update(id string, user *models.User) error {
	if err := s.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) Delete(id string) error {
	if err := s.db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (s *store) AddAvailability(slots []models.UserAvailability) error {
	for _, slot := range slots {
		if err := s.db.Create(&slot).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *store) UpdateAvailability(slotID string, slot models.UserAvailability) error {
	if err := s.db.Model(&models.UserAvailability{}).Where("id = ?", slotID).Updates(slot).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (s *store) DeleteAvailability(slotID string) error {
	if err := s.db.Where("id = ?", slotID).Delete(&models.UserAvailability{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (s *store) GetAvailability(userID string) ([]models.UserAvailability, error) {
	var availabilities []models.UserAvailability
	if err := s.db.Where("user_id = ? and event_id IS NULL", userID).Find(&availabilities).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return availabilities, nil
}
