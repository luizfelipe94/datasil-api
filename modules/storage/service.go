package storage

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/luizfelipe94/datasil/db"
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

func (s *Service) ListFiles(page int) (*db.ResultSet[models.File], error) {
	rowsCount, err := s.db.Query(`
		SELECT COUNT(*) AS count 
		FROM storage_files 
		WHERE deletedAt IS NULL
	`)
	if err != nil {
		return nil, err
	}
	pageCount := db.CountRows(rowsCount)
	limit := 10
	offset := limit * (page - 1)
	sql := `
		SELECT * 
		FROM storage_files 
		WHERE deletedAt IS NULL 
		ORDER BY createdAt DESC
		OFFSET $1
		LIMIT $2
	`
	rows, err := s.db.Query(sql, offset, limit)
	if err != nil {
		return nil, err
	}
	files := make([]models.File, 0)
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
			files = append(files, *file)
		}
	}
	res := db.NewResultSet[models.File](files, page, limit, pageCount)
	return &res, nil
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

func (s *Service) GetStats() (count int, size float64, average float64, err error) {
	rows, err := s.db.Query(`
		SELECT count(*) as count, sum(size) as size, avg(size) as average 
		FROM storage_files
	`)
	if err != nil {
		return 0, 0, 0, err
	}
	for rows.Next() {
		rows.Scan(&count, &size, &average)
	}
	return count, size, average, nil
}
