package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: config: cleanenv
	// TODO: logger: slog
	// TODO: storage: sqlite
	// TODO: router: chi
	// TODO: server: run server
}
