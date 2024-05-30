package models

import (
	"time"
)

type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Extension   string    `json:"extension"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType" validate:"required"`
	CraetedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"-"`
}
