package main

import (
	"fmt"
	"gliphtones/database"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config
var myCache *cache.Cache

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	database.Init()
	myCache = cache.New(10*time.Second, 1*time.Hour)

	e := echo.New()

	e.Static("/static", "static")
	setupRouter(e)

	port := fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))
	e.Logger.Fatal(e.Start(port))
}
