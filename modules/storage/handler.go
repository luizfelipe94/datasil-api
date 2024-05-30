package storage

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/luizfelipe94/datasil/configs"
	"github.com/luizfelipe94/datasil/infra"
	"github.com/luizfelipe94/datasil/modules/storage/dto"
	"github.com/luizfelipe94/datasil/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Handler struct {
	service *Service
	s3      *infra.S3
}

func NewStorageHandler(db *sql.DB) *Handler {
	minioClient, err := minio.New(configs.Envs.Aws.S3.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(configs.Envs.Aws.AccessKey, configs.Envs.Aws.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}
	s3 := infra.NewS3(minioClient)
	return &Handler{
		service: NewService(db),
		s3:      s3,
	}
}

func (r *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("GET /storage/files", r.handleListFiles)
	router.HandleFunc("POST /storage/files", r.handleUploadFile)
}

func (h *Handler) handleListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.service.ListFiles()
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.ResponseOk(w, http.StatusOK, files)
}

func (h *Handler) handleUploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, _ := r.FormFile("file")
	defer file.Close()
	// path := r.FormValue("path")
	tmpPath := filepath.Join("/tmp/", header.Filename)
	dst, _ := os.Create(tmpPath)
	defer dst.Close()
	io.Copy(dst, file)

	err := h.s3.UploadFile("datasil-storage", header.Filename, tmpPath, header.Size, header.Header.Get("Content-Type"), nil)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	dto := dto.CreateFileDto{
		Name:        utils.GetFileName(header.Filename),
		Extension:   utils.GetFileExtension(header.Filename),
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
	}
	if err := h.service.UploadFile(dto); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.ResponseOk(w, http.StatusCreated, map[string]string{"message": "File uploaded successfully"})
}
