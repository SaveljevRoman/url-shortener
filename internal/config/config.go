package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Storage    Storage    `yaml:"storage" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server" env-required:"true"`
}

type Storage struct {
	Address  string `yaml:"address" env-default:"localhost:54321"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"root"`
}

type HttpServer struct {
	Addr        string        `yaml:"addr" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Не нашел путь к конфиг файлу")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Нет конфиг файла по пути: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Не могу прочитать файл конфига: %s", err)
	}

	return &cfg
}
