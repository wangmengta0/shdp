package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig 全局配置结构体
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type MySQLConfig struct {
	DSN          string `mapstructure:"dsn"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
type RabbitMQConfig struct {
	URL       string `mapstructure:"url"`
	QueueName string `mapstructure:"queue_name"`
}

var Conf *AppConfig

// InitConfig 初始化并加载配置文件
func InitConfig() {
	viper.SetConfigFile("config/application.yaml") // 指定配置文件路径
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	Conf = &AppConfig{}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	log.Println("配置文件加载成功！")
}
