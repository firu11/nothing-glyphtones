package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"

	"glyphtones/database"
	"glyphtones/templates/components"
	"glyphtones/templates/views"

	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"
	godiacritics "gopkg.in/Regis24GmbH/go-diacritics.v2"
)

func (s *Server) googleLogin(c *echo.Context) error {
	url := s.googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) googleCallback(c *echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return Render(c, views.OtherErrorView(http.StatusBadRequest, errors.New("Bad request")))
	}

	token, err := s.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println(err)
		return Render(c, views.OtherErrorView(http.StatusInternalServerError, errors.New("Failed to exchange token")))
	}

	client := s.googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return Render(c, views.OtherErrorView(http.StatusInternalServerError, errors.New("Failed to fetch user info")))
	}
	defer resp.Body.Close()

	var authorInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&authorInfo); err != nil {
		return Render(c, views.OtherErrorView(http.StatusInternalServerError, errors.New("Failed to decode author info")))
	}

	name := authorInfo["name"].(string)
	name = godiacritics.Normalize(name)
	name = strings.Trim(name, " ")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	if len(name) > 30 {
		name = name[0:30]
	}
	if !authorNameR.MatchString(name) {
		name = fmt.Sprintf("author%d", rand.IntN(10000))
	}

	authorID, err := database.CreateAuthor(name, authorInfo["email"].(string))
	if err != nil {
		if strings.Contains(err.Error(), "unique_name") {
			authorID, _ = database.CreateAuthor(fmt.Sprintf("%s%d", name, rand.IntN(10000)), authorInfo["email"].(string))
		} else {
			return Render(c, views.OtherErrorView(http.StatusInternalServerError, err))
		}
	}

	if err := s.auth.WriteAuthCookie(c, authorID); err != nil {
		return Render(c, views.OtherErrorView(http.StatusInternalServerError, err))
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/me")
}

func (s *Server) logout(c *echo.Context) error {
	s.auth.RemoveAuthCookie(c)
	c.Response().Header().Set("HX-Redirect", "/")
	return Render(c, components.Header(false))
}
