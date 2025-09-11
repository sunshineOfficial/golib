package authorize

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

type Authorize struct {
	UserId  int    `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	Iat     int    `json:"iat"`
	Exp     int    `json:"exp"`
}

func GetToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func Parse(token string) (auth Authorize, err error) {
	token, ok := strings.CutPrefix(token, "Bearer ")
	if !ok {
		return Authorize{}, errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return Authorize{}, fmt.Errorf("could not parse token: %w", err)
	}

	return Authorize{
		UserId:  int(claims["user_id"].(float64)),
		Email:   claims["email"].(string),
		IsAdmin: claims["is_admin"].(bool),
		Iat:     int(claims["iat"].(float64)),
		Exp:     int(claims["exp"].(float64)),
	}, nil
}

func NewAccessToken(userId int, email string, isAdmin bool, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId
	claims["email"] = email
	claims["is_admin"] = isAdmin
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
