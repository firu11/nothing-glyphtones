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
)

var myCache *cache.Cache

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	database.Init()
	myCache = cache.New(10*time.Second, 1*time.Hour)

	e := echo.New()
	e.Static("/static", "static")
	setupRouter(e)

	port := fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))
	e.Logger.Fatal(e.Start(port))
}
