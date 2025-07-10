package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID           `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string              `gorm:"column:name;not null" json:"name"`
	Email          string              `gorm:"column:email;unique;not null" json:"email"`
	Availabilities []*UserAvailability `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserAvailability struct {
	ID      uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UserID  uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	EventID *uuid.UUID `gorm:"column:event_id;type:uuid"`
	Slot
}

type Slot struct {
	StartTime time.Time `gorm:"column:start_time;not null" json:"start_time"`
	EndTime   time.Time `gorm:"column:end_time;not null" json:"end_time"`
}
