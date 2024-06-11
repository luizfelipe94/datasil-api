package models

import (
	"time"
)

type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Extension   string    `json:"extension"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	CraetedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
	CompanyId   string    `json:"-"`
	Path        string    `json:"path"`
	IsFolder    bool      `json:"isFolder"`
	Depth       int       `json:"depth"`
}
