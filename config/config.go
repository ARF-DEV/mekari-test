package config

type Config struct {
	JWT_SECRET          string `mapstructure:"JWT_SECRET"`
	DB_MASTER           string `mapstructure:"DB_MASTER"`
	PAYMENT_GATEWAY_URL string `mapstructure:"PAYMENT_GATEWAY_URL"`
}
