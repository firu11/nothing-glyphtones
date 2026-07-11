package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"glyphtones/config"
	"glyphtones/database"
	"glyphtones/server"
	"glyphtones/utils"

	"github.com/labstack/echo/v5"
)

func configureLogging(production bool) {
	var handler slog.Handler
	if production {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		handler = slog.NewTextHandler(os.Stderr, nil)
	}
	slog.SetDefault(slog.New(handler))
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}
	configureLogging(cfg.Production)

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

	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%s", cfg.ListenPort),
		GracefulTimeout: 10 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
