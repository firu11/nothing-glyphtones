package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

func (s *Server) registerRoutes(e *echo.Echo) {
	e.RouteNotFound("/*", s.notFound)

	e.GET("/health", s.healthcheck)

	e.GET("/", s.index, fullPageGzipMiddleware())
	e.GET("/me", s.me)
	e.GET("/author/:name", s.author, fullPageGzipMiddleware())
	e.GET("/rename-author", s.authorRenameView)
	e.POST("/rename-author", s.authorRename)
	e.GET("/upload", s.uploadView, fullPageGzipMiddleware())
	e.PUT("/upload", s.uploadFile)
	e.POST("/vote/:displayID", s.vote)
	e.POST("/download/:displayID", s.downloadRingtone)
	e.GET("/rename/:displayID", s.renameView)
	e.POST("/rename/:displayID", s.rename)
	e.POST("/admin/ringtone/:displayID", s.updateRingtoneMetadata)
	e.POST("/delete-ringtone/:displayID", s.deleteRingtone)
	e.GET("/g/:displayID", s.detailRingtone, fullPageGzipMiddleware())
	e.GET("/guide", s.guide, fullPageGzipMiddleware())
	e.GET("/dmca", s.dmca, fullPageGzipMiddleware())
	e.GET("/google-login", s.googleLogin)
	e.GET("/google-callback", s.googleCallback)
	e.POST("/logout", s.logout)
}

func (s *Server) healthcheck(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	if err := s.store.PingContext(ctx); err != nil {
		slog.Error("database healthcheck failed", "method", c.Request().Method, "path", c.Request().URL.Path, "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
