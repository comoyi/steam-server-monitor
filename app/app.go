package app

import (
	"github.com/comoyi/steam-server-monitor/api"
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/comoyi/steam-server-monitor/config"
)

func Start() {
	initApp()
	if config.Conf.EnableApi {
		go func() {
			api.Start()
		}()
	}
	client.Start()
}

func initApp() {
	config.LoadConfig()
	_ = config.SaveConfig()
}
