package main

import (
	"fmt"
	"glyphtones/database"
	"glyphtones/utils"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

const maxRingtoneSize = 2 * 1024 * 1024 // 2MB

var ringtoneNameR regexp.Regexp = *regexp.MustCompile("^[ -~]{2,30}$")
var authorNameR regexp.Regexp = *regexp.MustCompile("^[ -~]{3,20}$")

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(utils.RingtonesDir); os.IsNotExist(err) {
		err := os.Mkdir(utils.RingtonesDir, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
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

	if os.Getenv("PRODUCTION") == "false" {
		e.Static("/static", "static")
		e.Static("/sounds", "sounds")
	}
	setupRouter(e)

	port := fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))
	e.Logger.Fatal(e.Start(port))
}
