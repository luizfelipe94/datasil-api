package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luizfelipe94/datasil/modules/auth"
	"github.com/luizfelipe94/datasil/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)
		if tokenString == "" {
			utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", claims["userId"])
		ctx = context.WithValue(ctx, "companyId", claims["companyId"])
		ctx = context.WithValue(ctx, "email", claims["email"])
		ctx = context.WithValue(ctx, "name", claims["name"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		return strings.Split(r.Header.Get("Authorization"), " ")[1]
	}
	return ""
}
