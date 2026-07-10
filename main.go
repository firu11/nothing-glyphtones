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
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(cfg.RingtonesDir, 0o755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(cfg.TemporaryDir, 0o755); err != nil {
		log.Fatal(err)
	}

	store, err := database.Open(ctx, cfg.DBConnectionString)
	if err != nil {
		log.Fatal("database init failed: ", err)
	}
	defer store.Close()

	auth := utils.NewAuth(cfg.TokenKey, cfg.Production)
	appServer := server.NewServer(cfg, store, auth)

	e, err := appServer.NewEcho()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Start(fmt.Sprintf(":%s", cfg.ListenPort)); err != nil {
		log.Fatal(err)
	}
}
