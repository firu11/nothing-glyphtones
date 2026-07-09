package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const CookieName = "GlyphtonesCookie"

const issuer = "glyphtones.firu.dev"

const tokenLifetime = 14 * 24 * time.Hour

type Auth struct {
	secure     bool
	privateKey []byte
}

type data struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func NewAuth(tokenKey string, secure bool) *Auth {
	return &Auth{
		secure:     secure,
		privateKey: []byte(tokenKey),
	}
}

func (a *Auth) generateToken(id int) (string, error) {
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
	return token.SignedString(a.privateKey)
}

func (a *Auth) validateToken(tokenString string) (bool, int, error) {
	data := data{}
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.privateKey, nil
	},
		jwt.WithLeeway(10*time.Second),
		jwt.WithIssuedAt(),
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

func (a *Auth) WriteAuthCookie(c *echo.Context, id int) error {
	jwt, err := a.generateToken(id)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     CookieName,
		Value:    jwt,
		Path:     "/",
		Expires:  time.Now().Add(tokenLifetime),
		HttpOnly: true,
		Secure:   a.secure,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)
	return nil
}

func (a *Auth) GetIDFromCookie(c *echo.Context) int {
	cookie, err := c.Cookie(CookieName)
	if err != nil {
		return 0
	}
	valid, id, err := a.validateToken(cookie.Value)
	if err != nil || !valid {
		a.RemoveAuthCookie(c)
		return 0
	}
	return id
}

func (a *Auth) RemoveAuthCookie(c *echo.Context) {
	cookie := http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.secure,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)
}
