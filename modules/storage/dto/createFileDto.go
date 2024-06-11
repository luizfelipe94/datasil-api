package dto

import (
	"fmt"
	"strings"
)

type CreateFileDto struct {
	Name        string `json:"name" validate:"required"`
	Extension   string `json:"extension" validate:"required"`
	Size        int64  `json:"size" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
	Path        string `json:"path" validate:"required"`
}

func (dto *CreateFileDto) GetDepth() int {
	fmt.Println(dto.Path)
	if dto.Path == "/" {
		return 0
	}
	fmt.Println(len(strings.Split(dto.Path, "/")) - 1)
	return len(strings.Split(dto.Path, "/")) - 1
}
