package client

import (
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/log"
)

var appName = "Steam服务器信息查看器"
var versionText = "1.0.2"

func Start() {
	log.Debugf("Client start\n")

	initUI()

	loadServers()

	go func() {
		run()
	}()

	w.ShowAndRun()
}

func loadServers() {
	for _, s := range config.Conf.Servers {
		server := NewServer(s.Ip, s.Port, s.Interval, s.Remark)
		serverContainer.AddServer(server)
	}
}

func run() {
	for _, server := range serverContainer.GetServers() {
		bind(server)
		server.Start()
	}
}
