package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var cfg Config

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.BindEnv("DB_MASTER")
	viper.BindEnv("PORT")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Log().Err(err).Msgf("error when reading config")
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
