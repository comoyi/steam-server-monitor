package app

import (
	"fmt"
	"github.com/comoyi/steam-server-monitor/config"
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
	g := gui.New()
	g.Run()
}
