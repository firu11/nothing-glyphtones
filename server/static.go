package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func ResolveDir(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return "", fmt.Errorf("directory %q not found: %w", dir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path %q is not a directory", dir)
	}

	return absDir, nil
}

func RegisterStaticRoutes(e *echo.Echo, staticDir string, soundsDir string) {
	if staticDir != "" {
		e.Static("/static", staticDir, middleware.Gzip(), cacheControlMiddleware("max-age=3600"))
	}
	if soundsDir != "" {
		e.Static("/sounds", soundsDir, cacheControlMiddleware("max-age=604800"))
	}
}

func cacheControlMiddleware(cacheControl string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, cacheControl)
			return next(c)
		}
	}
}

func dynamicNoCacheMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		path := c.Request().URL.Path
		if !strings.HasPrefix(path, "/static") && !strings.HasPrefix(path, "/sounds") {
			c.Response().Header().Set(echo.HeaderCacheControl, "no-store, no-cache, must-revalidate, max-age=0")
			c.Response().Header().Set("Pragma", "no-cache")
		}
		return next(c)
	}
}
