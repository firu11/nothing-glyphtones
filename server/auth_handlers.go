package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"glyphtones/templates/components"
	"glyphtones/templates/views"

	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"
	godiacritics "gopkg.in/Regis24GmbH/go-diacritics.v2"
)

type googleUserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Server) googleLogin(c *echo.Context) error {
	state, err := randomHex(16)
	if err != nil {
		return internalError(c, "generate OAuth state", err, true)
	}

	c.SetCookie(&http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.cfg.Production,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(10 * time.Minute),
	})

	url := s.googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) googleCallback(c *echo.Context) error {
	stateCookie, err := c.Cookie(oauthStateCookieName)
	if err != nil || stateCookie.Value == "" || c.QueryParam("state") != stateCookie.Value {
		return Render(c, views.OtherErrorView(http.StatusBadRequest, errors.New("Invalid login state")))
	}
	s.clearOAuthStateCookie(c)

	code := c.QueryParam("code")
	if code == "" {
		return Render(c, views.OtherErrorView(http.StatusBadRequest, errors.New("Bad request")))
	}

	ctx := c.Request().Context()
	token, err := s.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return internalError(c, "exchange Google OAuth token", err, true)
	}

	client := s.googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return internalError(c, "fetch Google user info", err, true)
	}
	defer resp.Body.Close()

	var authorInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&authorInfo); err != nil {
		return internalError(c, "decode Google user info", err, true)
	}
	if authorInfo.Name == "" || authorInfo.Email == "" {
		return internalError(c, "validate Google user info", errors.New("missing required profile fields"), true)
	}

	name := normalizeAuthorName(authorInfo.Name)
	authorID, err := s.store.CreateAuthor(ctx, name, authorInfo.Email)
	if err != nil {
		if strings.Contains(err.Error(), "unique_name") {
			authorID, err = s.store.CreateAuthor(ctx, fmt.Sprintf("%s%d", name, rand.IntN(10000)), authorInfo.Email)
		}
		if err != nil {
			return internalError(c, "create author", err, true)
		}
	}

	if err := s.auth.WriteAuthCookie(c, authorID); err != nil {
		return internalError(c, "write auth cookie", err, true)
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/me")
}

func (s *Server) logout(c *echo.Context) error {
	s.auth.RemoveAuthCookie(c)
	c.Response().Header().Set("HX-Redirect", "/")
	return Render(c, components.Header(false))
}

func normalizeAuthorName(name string) string {
	name = godiacritics.Normalize(name)
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	if len(name) > 30 {
		name = name[:30]
	}
	if !authorNameR.MatchString(name) {
		name = fmt.Sprintf("author%d", rand.IntN(10000))
	}
	return name
}

func (s *Server) clearOAuthStateCookie(c *echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     oauthStateCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.cfg.Production,
		SameSite: http.SameSiteLaxMode,
	})
}
