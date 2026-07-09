package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"glyphtones/config"
	"glyphtones/database"
	"glyphtones/server"
	"glyphtones/utils"
)

func main() {
	cfg, err := config.Load(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(utils.RingtonesDir, 0o755); err != nil {
		log.Panic(err)
	}

	if err := os.MkdirAll(utils.TemporaryDir, 0o755); err != nil {
		log.Panic(err)
	}

	if err := database.Init(cfg.DBConnectionString); err != nil {
		log.Fatal("database init failed: ", err)
	}

	auth := utils.NewAuth(cfg.TokenKey, cfg.Production)
	appServer := server.NewServer(cfg, auth)

	e, err := appServer.NewEcho()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Start(fmt.Sprintf(":%s", cfg.ListenPort)); err != nil {
		log.Fatal(err)
	}
}
