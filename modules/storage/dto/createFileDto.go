package dto

type CreateFileDto struct {
	Name        string `json:"name" validate:"required"`
	Extension   string `json:"extension" validate:"required"`
	Size        int64  `json:"size" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
}
