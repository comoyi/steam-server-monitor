package app

import (
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/spf13/viper"
)

func Start() {
	initApp()
	client.Start()
}

func initApp() {
	initConfig()
}

func initConfig() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.steam-server-monitor")
	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("Read config failed, err: %v\n", err)
		return
	}

	err = viper.Unmarshal(&config.Conf)
	if err != nil {
		log.Errorf("Unmarshal config failed, err: %v\n", err)
		return
	}
	log.Debugf("config: %+v\n", config.Conf)
}
