package config

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const DefaultConfigPath = "config/local.yaml"

var configPathFlag = flag.String("config", "", "path to YAML config file")

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Environment string `yaml:"app_env" env:"APP_ENV" env-default:"development"`
}

type ServerConfig struct {
	Host    string `yaml:"server_host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port    int    `yaml:"server_port" env:"SERVER_PORT" env-default:"8082"`
	Address string
}

type DatabaseConfig struct {
	URL string `yaml:"database_url" env:"DATABASE_URL"`
}

type JWTConfig struct {
	Secret             string `yaml:"jwt_secret" env:"JWT_SECRET"`
	AccessTokenTTL     string `yaml:"jwt_access_token_ttl" env:"JWT_ACCESS_TOKEN_TTL" env-default:"24h"`
	RefreshTokenTTL    string `yaml:"jwt_refresh_token_ttl" env:"JWT_REFRESH_TOKEN_TTL" env-default:"7d"`
	RefreshTokenSecret string `yaml:"jwt_refresh_token_secret" env:"JWT_REFRESH_TOKEN_SECRET"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"cors_allowed_origins" env:"CORS_ALLOWED_ORIGINS" env-separator:","`
	AllowedMethods   []string `yaml:"cors_allowed_methods" env:"CORS_ALLOWED_METHODS" env-separator:","`
	AllowedHeaders   []string `yaml:"cors_allowed_headers" env:"CORS_ALLOWED_HEADERS" env-separator:","`
	AllowCredentials bool     `yaml:"cors_allow_credentials" env:"CORS_ALLOW_CREDENTIALS" env-default:"false"`
}

func Load(configPath string) (Config, error) {
	var cfg Config

	yamlPath := configPath
	if yamlPath == "" {
		yamlPath = DefaultConfigPath
	}

	if err := cleanenv.ReadConfig(yamlPath, &cfg); err != nil {
		return cfg, fmt.Errorf("read %s: %w", yamlPath, err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, fmt.Errorf("read env overrides: %w", err)
	}
	cfg.Server.Address = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return cfg, nil
}

func LoadFromFlags() (Config, error) {
	flag.Parse()
	return Load(*configPathFlag)
}
