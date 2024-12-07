package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Static("/static", "static")
	setupRouter(e)

	port := fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))
	e.Logger.Fatal(e.Start(port))
}
