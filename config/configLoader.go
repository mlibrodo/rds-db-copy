package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var defaultConfig *Config

//go:embed default_config.toml
var defaultConfigFile []byte

func GetConfig() *Config {
	return defaultConfig
}

func init() {
	defaultConfig = readConfig()
}

func readConfig() *Config {
	var c Config
	viper.SetConfigType("toml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfigFile)); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}

	return &c
}
