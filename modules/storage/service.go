package storage

import (
	"database/sql"
	"fmt"

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
	query := fmt.Sprintf("SELECT * FROM files WHERE deletedAt IS NULL LIMIT 10")
	rows, err := s.db.Query(query)
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
