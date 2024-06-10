package storage

import (
	"context"
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

func (s *Service) ListFiles(ctx context.Context, page int) (*db.ResultSet[models.File], error) {
	companyId := ctx.Value("companyId")
	rowsCount, err := s.db.Query(`
		SELECT COUNT(*) AS count 
		FROM storage_files 
		WHERE deletedAt IS NULL
		AND companyId = $1
	`, companyId)
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
		AND companyId = $1
		ORDER BY createdAt DESC
		OFFSET $2
		LIMIT $3
	`
	rows, err := s.db.Query(sql, companyId, offset, limit)
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
			&file.CompanyId,
			&file.Path,
		)
		if err != nil {
			files = append(files, *file)
		}
	}
	res := db.NewResultSet[models.File](files, page, limit, pageCount)
	return &res, nil
}

func (s *Service) UploadFile(ctx context.Context, dto dto.CreateFileDto) error {
	id := uuid.New().String()
	companyId := ctx.Value("companyId")
	err := s.db.QueryRow(
		"INSERT INTO storage_files (id, name, extension, size, contentType, companyId) VALUES ($1, $2, $3, $4, $5, $6)",
		id, dto.Name, dto.Extension, dto.Size, dto.ContentType, companyId,
	)
	if err != nil {
		log.Println(err.Err())
		return err.Err()
	}
	return nil
}

func (s *Service) GetStats(ctx context.Context) (count int, size float64, average float64, err error) {
	companyId := ctx.Value("companyId")
	rows, err := s.db.Query(`
		SELECT count(*) as count, sum(size) as size, avg(size) as average 
		FROM storage_files
		WHERE deletedAt IS NULL
		AND companyId = $1
	`, companyId)
	if err != nil {
		return 0, 0, 0, err
	}
	for rows.Next() {
		rows.Scan(&count, &size, &average)
	}
	return count, size, average, nil
}
