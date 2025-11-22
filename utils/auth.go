package utils

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	privateKey        []byte        = []byte(os.Getenv("TOKEN_KEY"))
	tokenTimeDuration time.Duration = time.Hour * 24 * 14 // 14 days
	CookieName        string        = "GlyphtonesCookie"
)

type data struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func generateToken(id int) (string, error) {
	data := data{
		ID: id,
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
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (any, error) {
		return privateKey, nil
	})
	if err != nil {
		return false, 0, err
	}
	return token.Valid, data.ID, err
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
		Path:    "/",
	}
	c.SetCookie(&cookie)
	return nil
}

func GetIDFromCookie(c echo.Context) int {
	cookie, err := c.Cookie(CookieName)
	if err != nil {
		return 0
	}
	valid, id, err := validateToken(cookie.Value)
	if err != nil {
		return 0
	}
	if !valid {
		RemoveAuthCookie(c)
		return 0
	}
	return id
}

func RemoveAuthCookie(c echo.Context) {
	cookie := http.Cookie{
		Name:   CookieName,
		Value:  "",
		MaxAge: -1,
	}
	c.SetCookie(&cookie)
}
