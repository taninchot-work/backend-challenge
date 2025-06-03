package config

import "github.com/spf13/viper"

var config *Config

func InitConfig(configFileName string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return config
}

func SetConfig(c *Config) {
	config = c
}
