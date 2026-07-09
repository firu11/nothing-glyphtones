package server

import (
	"context"
	"net/http"
	"time"

	"glyphtones/database"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerRoutes(e *echo.Echo) {
	e.RouteNotFound("/*", s.notFound)

	e.GET("/health", s.healthcheck)

	e.GET("/", s.index)
	e.GET("/me", s.me)
	e.GET("/author/:name", s.author)
	e.GET("/rename-author", s.authorRenameView)
	e.POST("/rename-author", s.authorRename)
	e.GET("/upload", s.uploadView)
	e.PUT("/upload", s.uploadFile)
	e.POST("/vote/:displayID", s.vote)
	e.POST("/download/:displayID", s.downloadRingtone)
	e.GET("/rename/:displayID", s.renameView)
	e.POST("/rename/:displayID", s.rename)
	e.POST("/delete-ringtone/:displayID", s.deleteRingtone)
	e.GET("/g/:displayID", s.detailRingtone)
	e.GET("/guide", s.guide)
	e.GET("/dmca", s.dmca)
	e.GET("/google-login", s.googleLogin)
	e.GET("/google-callback", s.googleCallback)
	e.POST("/logout", s.logout)
}

func (s *Server) healthcheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	if database.DB.PingContext(ctx) != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
