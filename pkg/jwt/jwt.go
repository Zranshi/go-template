package jwt

import (
	"errors"
	"go-template/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	Email     string `json:"email"`
	Id        uint   `json:"id"`
	ExpiresAt int64  `json:"exp,omitempty"`
	jwt.StandardClaims
}

func GenerateToken(email string, uid uint, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Email: email,
		Id:    uid,
		ExpiresAt: time.Now().Add(
			config.Conf.GetDuration("host.tokenExpireDuration"),
		).Unix(),
	})
	return token.SignedString([]byte(key))
}

func ParseToken(tokenString string, key string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
