package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTIssuer       string
	AccessSecret    string
	RefreshSecret   string
	AccessMinutes   int
	RefreshDays     int
}

func Load() Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := Config{
		AppName: viper.GetString("APP_NAME"),
		AppPort: viper.GetString("APP_PORT"),

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),

		JWTIssuer:     viper.GetString("JWT_ISSUER"),
		AccessSecret:  viper.GetString("JWT_ACCESS_SECRET"),
		RefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
	}

	cfg.AccessMinutes = mustInt(viper.GetString("ACCESS_TOKEN_MINUTES"), 15)
	cfg.RefreshDays = mustInt(viper.GetString("REFRESH_TOKEN_DAYS"), 7)

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.AccessSecret == "" || cfg.RefreshSecret == "" {
		log.Fatal("JWT secrets are required")
	}
	return cfg
}

func mustInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
