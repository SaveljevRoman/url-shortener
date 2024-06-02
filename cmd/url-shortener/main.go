package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DbName,
		cfg.Storage.User,
		cfg.Storage.Password,
	)
	if err != nil {
		log.Error("ошибка подключения к хранилищу: %w", err)
	}
	fmt.Println(storage)

	// TODO: router: chi
	// TODO: server: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
