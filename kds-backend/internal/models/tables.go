package models

import (
	"time"
)

type Table struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Data        string    `json:"data"`
    CreatedBy   int       `json:"created_by"`
    IsPublic    bool      `json:"is_public"`
    CreatedAt   time.Time `json:"created_at"`
}

type TableRequest struct {
    Name        string `json:"name" binding:"required"`
    Description string `json:"description" binding:"required"`
    Data        string `json:"data" binding:"required"`
    IsPublic    bool   `json:"is_public"`
}