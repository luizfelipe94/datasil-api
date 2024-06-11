package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"

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

func (s *Service) ListFiles(ctx context.Context, page int, depth int) (*db.ResultSet[models.File], error) {
	companyId := ctx.Value("companyId")
	rowsCount, err := s.db.Query(`
		SELECT COUNT(*) AS count 
		FROM storage_files 
		WHERE deletedAt IS NULL
		AND companyId = $1
		AND depth = $2
	`, companyId, depth)
	if err != nil {
		return nil, err
	}
	pageCount := db.CountRows(rowsCount)
	limit := 10
	offset := limit * (page - 1)
	query := `
		SELECT id, name, extension, size, contentType, createdAt, path, isFolder, depth
		FROM storage_files 
		WHERE deletedAt IS NULL 
		AND companyId = $1
		AND depth = $2
		ORDER BY createdAt DESC
		OFFSET $3
		LIMIT $4
	`
	rows, err := s.db.QueryContext(ctx, query, companyId, depth, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]models.File, 0)
	for rows.Next() {
		file := new(models.File)
		var contentType sql.NullString
		var size sql.NullInt64
		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Extension,
			&size,
			&contentType,
			&file.CraetedAt,
			&file.Path,
			&file.IsFolder,
			&file.Depth,
		)
		if err != nil {
			log.Println("error:", err)
			return nil, err
		}
		file.ContentType = contentType.String
		file.Size = size.Int64
		files = append(files, *file)
	}
	res := db.NewResultSet[models.File](files, page, limit, pageCount)
	return &res, nil
}

func (s *Service) UploadFile(ctx context.Context, dto dto.CreateFileDto) error {
	id := uuid.New().String()
	companyId := ctx.Value("companyId")
	err := s.db.QueryRow(
		"INSERT INTO storage_files (id, name, extension, size, contentType, companyId, path, depth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, dto.Name, dto.Extension, dto.Size, dto.ContentType, companyId, dto.Path, dto.GetDepth(),
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
		AND isFolder = false
	`, companyId)
	if err != nil {
		return 0, 0, 0, err
	}
	for rows.Next() {
		rows.Scan(&count, &size, &average)
	}
	return count, size, average, nil
}

func (s *Service) CreateFolder(ctx context.Context, dto dto.CreateFolderDto) error {
	companyId := ctx.Value("companyId")
	basePath := ""
	depth := 0
	for _, path := range strings.Split(dto.Name, "/") {
		id := uuid.New().String()
		basePath += "/" + path
		_, err := s.db.ExecContext(
			ctx,
			"INSERT INTO storage_files (id, name, extension, companyId, isFolder, path, depth) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			id, path, "folder", companyId, true, basePath, depth,
		)
		if err != nil {
			log.Println("error:", err)
			return err
		}
		depth++
	}
	return nil
}
