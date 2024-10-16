package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/comoyi/steam-server-monitor/data"
)

func New() *Gui {
	gui := &Gui{}
	return gui
}

type Gui struct {
	Data       *data.Data
	App        fyne.App
	MainWindow *MainWindow
}

func (g *Gui) Run() {
	err := g.initApp()
	if err != nil {
		return
	}

	err = g.initMainWindow()
	if err != nil {
		return
	}

	g.MainWindow.Window.Show()
	g.App.Run()
}

func (g *Gui) initApp() error {
	a := app.NewWithID("com.comoyi.steamservermonitor")
	g.App = a
	return nil
}

func (g *Gui) initMainWindow() error {
	mainWindow := NewMainWindow()
	mainWindow.App = g.App
	mainWindow.Data = g.Data
	g.MainWindow = mainWindow
	err := mainWindow.Init()
	if err != nil {
		return err
	}
	return nil
}
