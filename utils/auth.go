package utils

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var privateKey []byte = []byte(os.Getenv("TOKEN_KEY"))
var tokenTimeDuration time.Duration = time.Hour * 24 * 14 // 14 days
var CookieName string = "GliphtonesCookie"

type data struct {
	Id int
	jwt.StandardClaims
}

func generateToken(id int) (string, error) {
	data := data{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTimeDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := token.SignedString(privateKey)
	return tokenString, err
}

func validateToken(tokenString string) (bool, int, error) {
	data := data{}
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return false, 0, err
	}
	return token.Valid, data.Id, err
}

func WriteAuthCookie(c echo.Context, id int) error {
	jwt, err := generateToken(id)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:    CookieName,
		Value:   jwt,
		Expires: time.Now().Add(tokenTimeDuration),
	}
	c.SetCookie(&cookie)
	return nil
}

func RemoveAuthCookie(c echo.Context) {
	cookie := http.Cookie{
		Name:   CookieName,
		Value:  "",
		MaxAge: -1,
	}
	c.SetCookie(&cookie)
}
