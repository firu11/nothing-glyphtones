package server

import (
	"regexp"

	"glyphtones/config"
	"glyphtones/utils"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	maxRingtoneSize      = 3 * 1024 * 1024 // 3MB
	lastSearchCookieName = "Glyphtones_last_search_options"
)

var (
	ringtoneNameR = *regexp.MustCompile("^[ -~]{2,30}$")
	authorNameR   = *regexp.MustCompile("^[a-z0-9_-]{3,20}$")
)

type Server struct {
	cfg               config.Config
	auth              *utils.Auth
	googleOauthConfig *oauth2.Config
}

func NewServer(cfg config.Config, auth *utils.Auth) *Server {
	return &Server{
		cfg:  cfg,
		auth: auth,
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  cfg.GoogleRedirectURL,
			ClientID:     cfg.GoogleID,
			ClientSecret: cfg.GoogleSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (s *Server) NewEcho() *echo.Echo {
	e := echo.New()

	if !s.cfg.Production {
		e.Static("/static", "static")
		e.Static("/sounds", "sounds")
	}

	s.registerRoutes(e)
	return e
}

func Render(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response())
}
