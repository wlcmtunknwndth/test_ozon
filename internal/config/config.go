package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	DB     Database `yaml:"db" env-required:"true"`
	Server Server   `yaml:"server" env-required:"true"`
	UseDB  bool     `yaml:"use_db" env-default:"false"`
}

type Database struct {
	DbUser  string `yaml:"db_user" env-required:"true"`
	DbPass  string `yaml:"db_pass" env-required:"true"`
	DbName  string `yaml:"db_name" env-required:"true"`
	SslMode string `yaml:"ssl_mode" env-default:"disable"`
	Port    string `yaml:"port" env-default:"5432"`
}

type Server struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout" env-default:"15s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

func MustLoad() *Config {
	const op = "config.MustLoad"

	path, ok := os.LookupEnv("config_path")
	if !ok || path == "" {
		slog.Error("couldn't find config path env:", slogAttr.SlogInfo(op, path))
		panic("couldn't find config path env")
	}

	if _, err := os.Stat(path); err != nil {
		slog.Error("couldn't find config path:", slogAttr.SlogErr(op, err))
		panic("couldn't find config path")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		slog.Error("couldn't read config:", slogAttr.SlogErr(op, err))
		panic("couldn't read config")
	}

	return &cfg
}
