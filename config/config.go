package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Production         bool   `env:"PRODUCTION,default=false"`
	ListenPort         string `env:"LISTEN_PORT,default=8080"`
	DBConnectionString string `env:"DB_CONNECTION_STRING,required"`
	GoogleID           string `env:"GOOGLE_ID,required"`
	GoogleSecret       string `env:"GOOGLE_SECRET,required"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL,required"`
	TokenKey           string `env:"TOKEN_KEY,required"`
	RingtonesDir       string `env:"RINGTONES_DIR,default=./sounds"`
	TemporaryDir       string `env:"TEMPORARY_DIR,default=./tmp"`
}

func Load(ctx context.Context) (Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
