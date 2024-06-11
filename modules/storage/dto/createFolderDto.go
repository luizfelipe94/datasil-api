package dto

type CreateFolderDto struct {
	Name string `json:"name" validate:"required"`
}
