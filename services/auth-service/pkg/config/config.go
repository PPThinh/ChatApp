package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	GRPCPort     string
	UserGRPCPort string
	JWTSecretKey string
}

func LoadConfig() *Config {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./pkg/config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Lỗi đọc file config: %v", err)
	}
	return &Config{
		GRPCPort:     viper.GetString("grpc.port"),
		UserGRPCPort: viper.GetString("user-grpc.port"),
		JWTSecretKey: viper.GetString("jwt.secret"),
	}
}
