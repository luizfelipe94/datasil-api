package storage

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/luizfelipe94/datasil/modules/storage/dto"
	"github.com/luizfelipe94/datasil/modules/storage/models"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) ListFiles() ([]*models.File, error) {
	rows, err := s.db.Query("SELECT * FROM storage_files WHERE deletedAt IS NULL LIMIT 10")
	if err != nil {
		return nil, err
	}
	files := make([]*models.File, 0)
	for rows.Next() {
		file := new(models.File)
		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Extension,
			&file.Size,
			&file.ContentType,
			&file.CraetedAt,
			&file.UpdatedAt,
			&file.DeletedAt,
		)
		if err != nil {
			files = append(files, file)
		}
	}
	return files, nil
}

func (s *Service) UploadFile(dto dto.CreateFileDto) error {
	id := uuid.New().String()
	err := s.db.QueryRow(
		"INSERT INTO storage_files (id, name, extension, size, contentType) VALUES ($1, $2, $3, $4, $5)",
		id, dto.Name, dto.Extension, dto.Size, dto.ContentType,
	)
	if err != nil {
		log.Println(err.Err())
		return nil
	}
	return nil
}
