package storage

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	api "github.com/luizfelipe94/datasil/cmd/api/middlewares"
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
	router.Handle("GET /storage/stats", api.AuthMiddleware(http.HandlerFunc(r.handleStats)))
	router.Handle("GET /storage/files", api.AuthMiddleware(http.HandlerFunc(r.handleListFiles)))
	router.Handle("POST /storage/files", api.AuthMiddleware(http.HandlerFunc(r.handleUploadFile)))
	router.Handle("POST /storage/folders", api.AuthMiddleware(http.HandlerFunc(r.handleCreateFolder)))
}

func (h *Handler) handleStats(w http.ResponseWriter, r *http.Request) {
	count, size, average, err := h.service.GetStats(r.Context())
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.ResponseOk(w, http.StatusCreated, map[string]any{
		"count":   count,
		"size":    size,
		"average": average,
	})
}

func (h *Handler) handleListFiles(w http.ResponseWriter, r *http.Request) {
	var page int = 1
	if r.URL.Query().Has("page") {
		t, _ := strconv.Atoi(r.URL.Query().Get("page"))
		page = t
	}
	depth := 0
	if r.URL.Query().Has("depth") {
		v, _ := strconv.Atoi(r.URL.Query().Get("depth"))
		depth = v
	}
	files, err := h.service.ListFiles(r.Context(), page, depth)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.ResponseOk(w, http.StatusOK, files)
}

func (h *Handler) handleUploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, _ := r.FormFile("file")
	path := r.FormValue("path")
	if path == "" {
		path = "/"
	}
	go func() {
		defer file.Close()
		tmpPath := filepath.Join("/tmp/", header.Filename)
		dst, _ := os.Create(tmpPath)
		defer dst.Close()
		io.Copy(dst, file)
		err := h.s3.UploadFile("datasil", header.Filename, tmpPath, header.Size, header.Header.Get("Content-Type"), nil)
		if err != nil {
			log.Println(err)
		}
	}()

	dto := dto.CreateFileDto{
		Name:        utils.GetFileName(header.Filename),
		Extension:   utils.GetFileExtension(header.Filename),
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Path:        path,
	}
	if err := h.service.UploadFile(r.Context(), dto); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.ResponseOk(w, http.StatusCreated, map[string]string{"message": "File uploaded successfully"})
}

func (h *Handler) handleCreateFolder(w http.ResponseWriter, r *http.Request) {
	var dto dto.CreateFolderDto
	if err := utils.ParseBody(r, &dto); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if err := utils.Validate.Struct(dto); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if err := h.service.CreateFolder(r.Context(), dto); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}
