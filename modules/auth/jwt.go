package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	users "github.com/luizfelipe94/datasil/modules/users/models"
)

func CreateJWT(secret []byte, user *users.User) (string, error) {
	expiration := time.Second * time.Duration(6000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.ID,
		"companyId": user.CompanyID,
		"email":     user.Email,
		"name":      user.Name,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
