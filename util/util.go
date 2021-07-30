package util

import (
	"github.com/spf13/viper"
)

func InitConfig(configName, configType, configPath string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig()
}
