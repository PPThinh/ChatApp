package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DSN  string
	Port string
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

	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.name"),
		viper.GetString("mysql.charset"),
		viper.GetBool("mysql.parseTime"),
		viper.GetString("mysql.loc"),
	)

	return &Config{
		DSN:  mysqlDSN,
		Port: viper.GetString("server.port"),
	}
}
