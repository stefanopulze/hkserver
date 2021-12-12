package configs

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

var Version = "0.0.1"

func Initialize() (*Config, error) {
	customConfigPath := flag.String("config", "", "--config <config-path")
	flag.Parse()

	if customConfigPath != nil && len(*customConfigPath) > 0 {
		viper.AddConfigPath(*customConfigPath)
	} else {
		viper.AddConfigPath("/etc/hkserver/")
		viper.AddConfigPath("$HOME/.hkserver")
		viper.AddConfigPath(".")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HK")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		return nil, fmt.Errorf("Fatal error config file: %w \n", err)
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &c, err
}
