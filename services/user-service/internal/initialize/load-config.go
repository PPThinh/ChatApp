package initialize

//
//import (
//	"fmt"
//	"log"
//
//	"github.com/spf13/viper"
//)
//
//type Config struct {
//	DSN      string
//	Port     string
//	LogLevel string
//	LogFile  string
//}
//
//func LoadConfig() {
//	viper := viper.New()
//	viper.SetConfigName("config")
//	viper.SetConfigType("yaml")
//	viper.AddConfigPath(".")
//
//	err := viper.ReadInConfig()
//	if err != nil {
//		log.Fatalf("Lỗi đọc file config: %v", err)
//	}
//
//	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
//		viper.GetString("mysql.user"),
//		viper.GetString("mysql.password"),
//		viper.GetString("mysql.host"),
//		viper.GetInt("mysql.port"),
//		viper.GetString("mysql.name"),
//		viper.GetString("mysql.charset"),
//		viper.GetBool("mysql.parseTime"),
//		viper.GetString("mysql.loc"),
//	)
//
//	return Config{
//		MySQLDSN: mysqlDSN,
//		Port:     viper.GetString("server.port"),
//		LogLevel: viper.GetString("log.level"),
//		LogFile:  viper.GetString("log.output"),
//	}
//}
