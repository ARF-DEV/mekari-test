package config

type Config struct {
	PORT      string `mapstructure:"PORT"`
	DB_MASTER string `mapstructure:"DB_MASTER"`
}
