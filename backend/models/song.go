package models

import (
	"time"

	"gorm.io/gorm"
)

type Song struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Title       string     `json:"title" gorm:"not null;index"`
	Artist      string     `json:"artist" gorm:"index"`
	Album       string     `json:"album" gorm:"index"`
	Genre       string     `json:"genre"`
	Year        int        `json:"year"`
	TrackNumber int        `json:"track_number"`
	Filename    string     `json:"filename" gorm:"not null"`
	FilePath    string     `json:"file_path" gorm:"not null"`
	FileSize    int64      `json:"file_size"`
	Duration    int        `json:"duration"` // seconds
	Bitrate     int        `json:"bitrate"`
	SampleRate  int        `json:"sample_rate"`
	MimeType    string     `json:"mime_type"`
	IsDeleted   bool       `json:"is_deleted" gorm:"default:false;index"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (s *Song) BeforeCreate(tx *gorm.DB) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = time.Now()
	}
	return nil
}

func (s *Song) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
