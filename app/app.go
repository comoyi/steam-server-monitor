package app

import (
	"fmt"
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/data"
	"github.com/comoyi/steam-server-monitor/gui"
)

func Run() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("load config failed, err: %v\n", err)
		return
	}
	fmt.Printf("load config success, conf: %+v\n", conf)
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
		fmt.Printf("config is nil\n")
		return
	}

	data := data.New()
	err := data.Init()
	if err != nil {
		fmt.Printf("init data failed, err: %v\n", err)
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
