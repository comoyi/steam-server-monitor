package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow() *MainWindow {
	return &MainWindow{}
}

type MainWindow struct {
	App     fyne.App
	Window  fyne.Window
	ToolBar *fyne.Container
}

func (w *MainWindow) Init() error {
	windowTitle := fmt.Sprintf("%s - v%s", "Steam Server Monitor", "2.0.1")
	window := w.App.NewWindow(windowTitle)
	w.Window = window
	window.Resize(fyne.NewSize(600, 400))

	c := container.NewVBox()

	c2 := container.NewVBox()
	toolBar, err := w.initToolbar()
	if err != nil {
		return err
	}
	c2.Add(toolBar)

	title := binding.NewString()
	err = title.Set("Title")
	if err != nil {
		fmt.Printf("set title failed, err: %v\n", err)
		return err
	}

	l := widget.NewLabelWithData(title)
	c2.Add(l)
	cScroll := container.NewVScroll(c2)
	cScroll.SetMinSize(fyne.NewSize(600, 300))

	c.Add(cScroll)
	w.Window.SetContent(c)

	return nil
}

func (w *MainWindow) initToolbar() (*fyne.Container, error) {
	bar := container.NewGridWithColumns(2)
	addBtn := widget.NewButton("+", func() {
		message := "TODO"
		content := container.NewVBox()
		messageLabel := widget.NewLabel(message)
		content.Add(messageLabel)
		dialog.NewCustom("Tip", "OK", content, w.Window).Show()
	})
	bar.Add(addBtn)
	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {

	})
	bar.Add(saveBtn)
	return bar, nil
}
