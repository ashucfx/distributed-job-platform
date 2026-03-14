package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	Port     string `mapstructure:"PORT"`
	DBUrl    string `mapstructure:"DATABASE_URL"`
	RedisUrl string `mapstructure:"REDIS_URL"`
}

var Cfg *Config

func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	viper.SetDefault("REDIS_URL", "localhost:6379")

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	Cfg = &config
	return Cfg
}
