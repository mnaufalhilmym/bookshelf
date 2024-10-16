package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	v := viper.New()

	v.SetConfigFile("config.yml")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	return v
}
