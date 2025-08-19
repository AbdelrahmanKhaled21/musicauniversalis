package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Name          string         `json:"name"`
	OAuthProvider string         `json:"oauth_provider" gorm:"column:oauth_provider"`
	OAuthID       string         `json:"oauth_id" gorm:"column:oauth_id"`
	StorageQuota  int64          `json:"storage_quota" gorm:"default:10737418240"` // 10GB default
	StorageUsed   int64          `json:"storage_used" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Songs     []Song     `json:"songs,omitempty" gorm:"foreignKey:UserID"`
	Playlists []Playlist `json:"playlists,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
