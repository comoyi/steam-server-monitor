package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`
}

func LoadConfig() (*Config, error) {
	var conf *Config

	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fixedExePath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(fixedExePath)

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(exeDir)
	v.AddConfigPath(filepath.Join(exeDir, "config"))

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	usedConfig := v.ConfigFileUsed()
	fmt.Printf("used config: %s\n", usedConfig)

	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
