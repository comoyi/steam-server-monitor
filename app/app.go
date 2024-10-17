package app

import (
	"fmt"
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/data"
	"github.com/comoyi/steam-server-monitor/gui"
	"github.com/comoyi/steam-server-monitor/log"
)

func Run() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("load config failed, err: %v\n", err)
		return
	}
	fmt.Printf("load config success, conf: %+v\n", conf)

	err = log.Init()
	if err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	_, err = log.SetLogLevelByName(conf.LogLevel)
	if err != nil {
		fmt.Printf("set log level failed, err: %v\n", err)
		return
	}

	log.Debugf("log level: %s", log.LogLevel())

	a := New(conf)
	a.Run()
}

func New(conf *config.Config) *App {
	app := &App{
		conf: conf,
	}
	return app
}

type App struct {
	conf *config.Config
}

func (a *App) Run() {
	if a.conf == nil {
		log.Errorf("config is nil")
		return
	}

	data := data.New()
	err := data.Init()
	if err != nil {
		log.Errorf("init data failed, err: %v", err)
		return
	}

	go func() {
		client := client.New()
		client.Data = data
		client.Run()
	}()

	g := gui.New()
	g.Data = data
	g.Run()
}
