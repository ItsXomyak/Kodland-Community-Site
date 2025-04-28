package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Email               string    `json:"email" db:"email"`
	PasswordHash        string    `json:"-" db:"password_hash"`
	DisplayName         string    `json:"displayName" db:"display_name"`
	IsActive            bool      `json:"isActive" db:"is_active"`
	Role                UserRole  `json:"role" db:"role"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
	FailedLoginAttempts int       `json:"failedLoginAttempts" db:"failed_login_attempts"`
	LastFailedLogin     time.Time `json:"lastFailedLogin" db:"last_failed_login"`
}

type UserProfile struct {
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Login     string    `json:"login" db:"login"`
	Bio       string    `json:"bio" db:"bio"`
	AvatarURL string    `json:"avatarUrl" db:"avatar_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ProfileComment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	ProfileID uuid.UUID `json:"profileId" db:"profile_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
