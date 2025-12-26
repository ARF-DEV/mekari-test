package config

type Config struct {
	JWT_SECRET string `mapstructure:"JWT_SECRET"`
	PORT       string `mapstructure:"PORT"`
	DB_MASTER  string `mapstructure:"DB_MASTER"`
}
