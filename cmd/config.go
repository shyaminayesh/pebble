package cmd

import (
	"fmt"
	"github.com/spf13/viper"
)

func Config() *viper.Viper {

	var v = viper.New()
	v.SetConfigName("pebble")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return v
}