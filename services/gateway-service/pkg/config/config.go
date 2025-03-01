package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port            string
	GrpcUserAddress string
	GrpcAuthAddress string
}

func LoadConfig() *Config {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./pkg/config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return &Config{
		Port:            viper.GetString("server.port"),
		GrpcUserAddress: viper.GetString("grpc.user.address"),
		GrpcAuthAddress: viper.GetString("grpc.auth.address"),
	}
}
