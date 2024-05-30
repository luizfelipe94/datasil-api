package api

import (
	"database/sql"
	"net/http"

	"github.com/luizfelipe94/datasil/modules/storage"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	storageHandler := storage.NewStorageHandler(s.db)
	storageHandler.Register(router)

	return http.ListenAndServe(s.addr, router)
}
