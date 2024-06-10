package auth

import (
	"database/sql"
	"net/http"

	"github.com/luizfelipe94/datasil/modules/auth/dto"
	"github.com/luizfelipe94/datasil/utils"
)

type Handler struct {
	service *Service
}

func NewSAuthHandler(db *sql.DB) *Handler {
	return &Handler{
		service: NewService(db),
	}
}

func (r *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("POST /auth/login", r.login)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var dto dto.LoginDTO
	if err := utils.ParseBody(r, &dto); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if err := utils.Validate.Struct(dto); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	user, err := h.service.GetUserByEmail(dto.Email)
	if err != nil {
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !ComparePasswords(user.Password, []byte(dto.Password)) {
		utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	token, err := CreateJWT([]byte("flamengo@2024"), user)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.ResponseOk(w, http.StatusOK, map[string]any{
		"token": token,
	})
}
