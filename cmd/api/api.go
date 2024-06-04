package api

import (
	"database/sql"
	"net/http"

	"github.com/luizfelipe94/datasil/modules/storage"
	"github.com/rs/cors"
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
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok!"))
	})

	storageHandler := storage.NewStorageHandler(s.db)
	storageHandler.Register(router)

	handler := cors.AllowAll().Handler(router)
	return http.ListenAndServe(s.addr, handler)
}
