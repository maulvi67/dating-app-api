package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string
	Gender       string
	IsPremium    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Swipe struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index;not null"`
	TargetUserID uint      `gorm:"not null"`
	Action       string    `gorm:"not null"` // "like" or "pass"
	SwipeDate    time.Time `gorm:"index;not null"`
	CreatedAt    time.Time
}

// Claims for JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
