package server

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"

	"glyphtones/config"
	"glyphtones/database"
	"glyphtones/templates/views"
	"glyphtones/utils"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	maxRingtoneSize      = 3 * 1024 * 1024 // 3MB
	lastSearchCookieName = "Glyphtones_last_search_options"
	oauthStateCookieName = "Glyphtones_oauth_state"
)

var (
	ringtoneNameR = *regexp.MustCompile("^[ -~]{2,30}$")
	authorNameR   = *regexp.MustCompile("^[a-z0-9_-]{3,20}$")
)

type Server struct {
	cfg               config.Config
	store             *database.Store
	auth              *utils.Auth
	googleOauthConfig *oauth2.Config
}

func NewServer(cfg config.Config, store *database.Store, auth *utils.Auth) *Server {
	return &Server{
		cfg:   cfg,
		store: store,
		auth:  auth,
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  cfg.GoogleRedirectURL,
			ClientID:     cfg.GoogleID,
			ClientSecret: cfg.GoogleSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (s *Server) NewEcho() (*echo.Echo, error) {
	e := echo.New()
	e.Logger = slog.Default()
	e.Use(errorLoggingMiddleware)
	e.Use(dynamicNoCacheMiddleware)

	staticDir, err := ResolveDir("static")
	if err != nil {
		if s.cfg.Production {
			return nil, err
		}
		slog.Warn("failed to resolve static directory", "error", err)
	}

	soundsDir, err := ResolveDir(s.cfg.RingtonesDir)
	if err != nil {
		if s.cfg.Production {
			return nil, err
		}
		slog.Warn("failed to resolve sounds directory", "error", err)
	}

	RegisterStaticRoutes(e, staticDir, soundsDir)
	s.registerRoutes(e)
	return e, nil
}

func Render(c *echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response())
}

func errorLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		err := next(c)
		if err != nil {
			slog.Error("handler error", "method", c.Request().Method, "path", c.Request().URL.Path, "error", err)
		}
		return err
	}
}

func internalError(c *echo.Context, operation string, err error, fullPage bool) error {
	slog.Error("internal request error", "method", c.Request().Method, "path", c.Request().URL.Path, "operation", operation, "error", err)
	publicErr := errors.New("Something went wrong")
	if fullPage {
		return Render(c, views.OtherErrorView(http.StatusInternalServerError, publicErr))
	}
	return Render(c, views.OtherError(http.StatusInternalServerError, publicErr))
}
