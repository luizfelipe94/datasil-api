package auth

import (
	"fmt"
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

func ValidateToken(tokenString string) (*jwt.Token, error) {
	var secretKey = []byte("flamengo@2024")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signature method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error when parsing the token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return token, nil
}
