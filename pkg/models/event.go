package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID                uuid.UUID   `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title             string      `gorm:"column:title;not null" json:"title"`
	Description       string      `gorm:"column:description;type:text" json:"description"`
	EstimatedDuration int         `gorm:"column:estimated_duration;type:int;not null" json:"estimated_duration"`
	EventSlots        []EventSlot `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"event_slots"`
	OrganizerID       *uuid.UUID  `gorm:"column:organizer_id;type:uuid" json:"organizer_id"`
	Organizer         *User       `gorm:"foreignKey:ID" json:"-"`
}

type EventSlot struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	EventID   *uuid.UUID `gorm:"column:event_id;type:uuid" json:"event_id"`
	StartTime time.Time  `gorm:"column:start_time;not null" json:"start_time"`
	EndTime   time.Time  `gorm:"column:end_time;not null" json:"end_time"`
}

type RecommendedSlot struct {
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	UserIDs        []string  `json:"user_ids"`         // users who can attend
	MissingUserIDs []string  `json:"missing_user_ids"` // users who can't
}
