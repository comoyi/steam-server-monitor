package app

import (
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/comoyi/steam-server-monitor/config"
)

func Start() {
	initApp()
	client.Start()
}

func initApp() {
	config.LoadConfig()
	config.SaveConfig()
}
