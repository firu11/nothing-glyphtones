package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	secure     bool   = os.Getenv("PRODUCTION") == "true"
	privateKey []byte = []byte(os.Getenv("TOKEN_KEY"))

	tokenLifetime time.Duration = 14 * 24 * time.Hour // 14 days
	CookieName                  = "GlyphtonesCookie"
	issuer                      = "glyphtones.firu.dev"
)

type data struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func generateToken(id int) (string, error) {
	claims := data{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{issuer},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKey)
}

func validateToken(tokenString string) (bool, int, error) {
	data := data{}
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return privateKey, nil
	},
		jwt.WithLeeway(10*time.Second), // allow small clock skew
		jwt.WithIssuedAt(),             // validate iat
		jwt.WithIssuer(issuer),
		jwt.WithAudience(issuer),
	)
	if err != nil {
		return false, 0, err
	}
	if !token.Valid {
		return false, 0, fmt.Errorf("token claims invalid (possibly wrong issuer/audience/clock)")
	}
	return token.Valid, data.ID, err
}

func WriteAuthCookie(c echo.Context, id int) error {
	jwt, err := generateToken(id)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     CookieName,
		Value:    jwt,
		Path:     "/",
		Expires:  time.Now().Add(tokenLifetime),
		HttpOnly: true,   // prevents JS access
		Secure:   secure, // set to false only in local dev (no HTTPS)
		SameSite: http.SameSiteLaxMode,
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
	if err != nil || !valid {
		RemoveAuthCookie(c)
		return 0
	}
	return id
}

func RemoveAuthCookie(c echo.Context) {
	cookie := http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)
}
