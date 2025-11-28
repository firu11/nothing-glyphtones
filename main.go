package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"glyphtones/database"
	"glyphtones/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

const maxRingtoneSize = 3 * 1024 * 1024 // 2MB

var (
	ringtoneNameR regexp.Regexp = *regexp.MustCompile("^[ -~]{2,30}$")
	authorNameR   regexp.Regexp = *regexp.MustCompile("^[a-z0-9_-]{3,20}$")
)

var LastSearchCookieName string = "Glyphtones_last_search_options"

func main() {
	if err := os.MkdirAll(utils.RingtonesDir, 0o755); err != nil {
		log.Panic(err)
	}

	if err := os.MkdirAll(utils.RingtonesDir, 0o755); err != nil {
		log.Panic(err)
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	database.Init()

	e := echo.New()

	if os.Getenv("PRODUCTION") == "false" || os.Getenv("PRODUCTION") == "" {
		e.Static("/static", "static")
		e.Static("/sounds", "sounds")
	}
	setupRouter(e)

	// TODO some env config
	port, ok := os.LookupEnv("LISTEN_PORT")
	if !ok {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(port))
}
