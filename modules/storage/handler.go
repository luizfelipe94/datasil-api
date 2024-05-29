package storage

import "net/http"

type Handler struct {
}

func NewStorageHandler() *Handler {
	return &Handler{}
}

func (r *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("POST /storage", r.handleListFiles)
}

func (h *Handler) handleListFiles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
