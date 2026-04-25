package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	Server ServerConfig `yaml:"server" mapstructure:"server"`
}

type ServerConfig struct {
	Port string `yaml:"port" mapstructure:"port"`
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 兼容：在 service 目录下运行（./config），或在仓库根目录运行（./service/config）
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./service/config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v（请确认从 service 目录运行，或保证存在 service/config/config.yaml）", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct:%v", err)
	}

}
