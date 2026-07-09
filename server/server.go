package server

import (
	"log"
	"regexp"

	"glyphtones/config"
	"glyphtones/utils"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

func (s *Server) NewEcho() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(dynamicNoCacheMiddleware)

	staticDir, err := ResolveDir("static")
	if err != nil {
		if s.cfg.Production {
			return nil, err
		}
		log.Println(err)
	}

	soundsDir, err := ResolveDir(utils.RingtonesDir)
	if err != nil {
		if s.cfg.Production {
			return nil, err
		}
		log.Println(err)
	}

	RegisterStaticRoutes(e, staticDir, soundsDir)
	s.registerRoutes(e)
	return e, nil
}

func Render(c *echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response())
}
