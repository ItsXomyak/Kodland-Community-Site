package models

import "time"

type User struct {
	ID          int            `json:"id"`
	Nickname    string         `json:"nickname"`
	AvatarURL   string         `json:"avatar_url"`
	Links       []string       `json:"links"`
	Achievements []string      `json:"achievements"`
	Roles       []string       `json:"roles"`
	Medals      map[string]int `json:"medals"`
	Likes       int            `json:"likes"`
	JoinedAt    time.Time      `json:"joined_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

type CreateUserRequest struct {
	Nickname    string         `json:"nickname" binding:"required"`
	AvatarURL   string         `json:"avatar_url" binding:"required"`
	Links       []string       `json:"links"`
	Achievements []string      `json:"achievements"`
	Roles       []string       `json:"roles"`
	Medals      map[string]int `json:"medals"`
	Likes       int            `json:"likes"`
	JoinedAt    time.Time      `json:"joined_at"`
}

type UpdateUserRequest struct {
	Nickname    string         `json:"nickname"`
	AvatarURL   string         `json:"avatar_url"`
	Links       []string       `json:"links"`
	Achievements []string      `json:"achievements"`
	Roles       []string       `json:"roles"`
	Medals      map[string]int `json:"medals"`
	Likes       *int           `json:"likes"`
	JoinedAt    *time.Time     `json:"joined_at"` 
}