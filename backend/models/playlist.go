package models

import (
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	User  User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Songs []PlaylistSong `json:"songs,omitempty" gorm:"foreignKey:PlaylistID"`
}

type PlaylistSong struct {
	PlaylistID uint      `json:"playlist_id" gorm:"primaryKey"`
	SongID     uint      `json:"song_id" gorm:"primaryKey"`
	Position   int       `json:"position" gorm:"not null"`
	AddedAt    time.Time `json:"added_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Playlist Playlist `json:"playlist,omitempty" gorm:"foreignKey:PlaylistID"`
	Song     Song     `json:"song,omitempty" gorm:"foreignKey:SongID"`
}

func (p *Playlist) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
	return nil
}

func (p *Playlist) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
