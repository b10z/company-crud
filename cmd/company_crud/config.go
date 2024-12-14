package main

import "github.com/spf13/viper"

type Config struct {
	Environment    string `mapstructure:"ENV"`
	Swagger        bool   `mapstructure:"SWAGGER"`
	Cors           bool   `mapstructure:"CORS"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	KafkaServer    string `mapstructure:"KAFKA_SERVER"`
	KafkaTopic     string `mapstructure:"KAFKA_TOPIC"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBName         string `mapstructure:"DB_NAME"`
	DBUsername     string `mapstructure:"DB_USERNAME"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	HTTPPort       int    `mapstructure:"HTTP_PORT"`
	TokenSignature string `mapstructure:"JWT_TOKEN_SIGNATURE"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
