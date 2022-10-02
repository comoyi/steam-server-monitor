package client

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	theme2 "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/comoyi/steam-server-monitor/theme"
	"github.com/comoyi/steam-server-monitor/util/dialogutil"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var w fyne.Window
var c *fyne.Container
var myApp fyne.App

var serverListPanel *fyne.Container

var serverListPanelScroll *container.Scroll

func initUI() {
	initMainWindow()
	initMenu()
}

func initMainWindow() {
	windowTitle := fmt.Sprintf("%s-v%s", appName, versionText)

	myApp = app.New()
	myApp.Settings().SetTheme(theme.CustomTheme)
	w = myApp.NewWindow(windowTitle)
	w.SetMaster()
	w.Resize(fyne.NewSize(400, 600))
	c = container.NewVBox()
	w.SetContent(c)

	bar := initToolBar()
	c.Add(bar)

	serverListPanel = container.NewVBox()
	serverListPanelScroll = container.NewVScroll(serverListPanel)
	serverListPanelScroll.SetMinSize(fyne.NewSize(400, 600))
	c.Add(serverListPanelScroll)

}

func initMenu() {
	addMenuItem := fyne.NewMenuItem("添加服务器", func() {
		showAddUI()
	})
	firstMenu := fyne.NewMenu("操作", addMenuItem)
	helpMenuItem := fyne.NewMenuItem("关于", func() {
		content := container.NewVBox()
		appInfo := widget.NewLabel(appName)
		content.Add(appInfo)
		versionInfo := widget.NewLabel(fmt.Sprintf("Version %v", versionText))
		content.Add(versionInfo)

		h := container.NewHBox()

		authorInfo := widget.NewLabel("Copyright © 2022 清新池塘")
		h.Add(authorInfo)
		linkInfo := widget.NewHyperlink(" ", nil)
		_ = linkInfo.SetURLFromString("https://github.com/comoyi/steam-server-monitor")
		h.Add(linkInfo)
		content.Add(h)
		dialog.NewCustom("关于", "关闭", content, w).Show()
	})
	helpMenu := fyne.NewMenu("帮助", helpMenuItem)
	mainMenu := fyne.NewMainMenu(firstMenu, helpMenu)
	w.SetMainMenu(mainMenu)
}

func initToolBar() *fyne.Container {
	cBar := container.NewGridWithColumns(2)

	addBtn := widget.NewButtonWithIcon("", theme2.ContentAddIcon(), func() {
		showAddUI()
	})
	cBar.Add(addBtn)

	var saveBtn *widget.Button
	saveText := "保存"
	saveBtn = widget.NewButtonWithIcon(saveText, theme2.DocumentSaveIcon(), func() {
		saveBtn.Disable()
		go func() {
			defer saveBtn.Enable()
			saveBtn.SetText("保存中...")
			log.Debugf("%+v\n", viper.AllSettings())
			err := config.SaveConfig()
			if err != nil {
				dialogutil.ShowInformation("提示", "保存失败", w)
				return
			}
			go func() {
				saveSuccessText := "保存成功"
				saveBtn.SetText(saveSuccessText)
				<-time.After(2 * time.Second)
				if saveBtn.Text == saveSuccessText {
					saveBtn.SetText(saveText)
				}
			}()
		}()
	})
	cBar.Add(saveBtn)

	return cBar
}

func showAddUI() {
	showServerFormUI(false, nil)
}

func showEditUI(server *Server) {
	showServerFormUI(true, server)
}

func showServerFormUI(isEdit bool, server *Server) {
	title := "添加服务器"
	if isEdit {
		if server == nil {
			return
		}
		title = "编辑服务器"
	}
	var serverFormWindow fyne.Window
	serverFormWindow = myApp.NewWindow(title)
	c := container.NewVBox()
	c2 := container.NewAdaptiveGrid(2)
	c3 := container.NewAdaptiveGrid(2)
	c4 := container.NewAdaptiveGrid(2)
	c5 := container.NewAdaptiveGrid(2)
	ipLabel := widget.NewLabel("IP")
	var ipEntry *widget.Entry
	ipEntry = widget.NewEntry()
	ipEntry.SetPlaceHolder("127.0.0.1")
	if isEdit {
		ipEntry.SetText(server.Ip)
	}

	portLabel := widget.NewLabel("端口")
	portHelpBtn := widget.NewButtonWithIcon("", theme2.HelpIcon(), func() {
		dialogutil.ShowInformation("", "信息查询端口，\n和主端口可能不同", serverFormWindow)
	})
	portBox := container.NewHBox()
	portBox.Add(portLabel)
	portBox.Add(portHelpBtn)
	var portEntry *widget.Entry
	portEntry = widget.NewEntry()
	portEntry.SetPlaceHolder("2457")
	if isEdit {
		portEntry.SetText(strconv.FormatInt(server.Port, 10))
	}
	intervalLabel := widget.NewLabel("刷新间隔（秒）")
	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder("10")
	intervalText := "10"
	if isEdit {
		intervalText = strconv.FormatInt(server.Interval, 10)
	}
	intervalEntry.Text = intervalText

	remarkLabel := widget.NewLabel("备注")
	var remarkEntry *widget.Entry
	remarkEntry = widget.NewEntry()
	if isEdit {
		remarkEntry.SetText(server.Remark)
	}

	btnText := "添加"
	if isEdit {
		btnText = "保存"
	}
	submitBtn := widget.NewButton(btnText, func() {
		ip := ipEntry.Text
		if ip == "" {
			dialogutil.ShowInformation("提示", "请输入IP", serverFormWindow)
			return
		}

		portVal := portEntry.Text
		if portVal == "" {
			dialogutil.ShowInformation("提示", "请输入端口", serverFormWindow)
			return
		}
		port, err := strconv.ParseInt(portVal, 10, 64)
		if err != nil {
			dialogutil.ShowInformation("提示", "请输入正确的端口", serverFormWindow)
			return
		}
		if port < 0 {
			dialogutil.ShowInformation("提示", "请输入正确的端口", serverFormWindow)
			return
		}

		intervalVal := intervalEntry.Text
		if intervalVal == "" {
			dialogutil.ShowInformation("提示", "请输入间隔", serverFormWindow)
			return
		}
		interval, err := strconv.ParseInt(intervalVal, 10, 64)
		if err != nil {
			dialogutil.ShowInformation("提示", "请输入正确的间隔", serverFormWindow)
			return
		}
		if interval <= 0 {
			dialogutil.ShowInformation("提示", "请输入合适的间隔", serverFormWindow)
			return
		}

		remark := remarkEntry.Text

		if isEdit {
			server.Ip = ip
			server.Port = port
			server.UpdateInterval(interval)
			server.Remark = remark
			refreshUI(server)
		} else {
			newServer := NewServer(ip, port, interval, remark)
			serverContainer.AddServer(newServer)
			bind(newServer)
			newServer.Start()
		}

		resetServerConfig()

		err = config.SaveConfig()
		if err != nil {
			dialogutil.ShowInformation("提示", "保存失败", w)
			return
		}

		serverFormWindow.Hide()
	})

	c2.Add(ipLabel)
	c2.Add(ipEntry)
	c3.Add(portBox)
	c3.Add(portEntry)
	c4.Add(intervalLabel)
	c4.Add(intervalEntry)
	c5.Add(remarkLabel)
	c5.Add(remarkEntry)
	c.Add(c2)
	c.Add(c3)
	c.Add(c4)
	c.Add(c5)
	c.Add(submitBtn)

	serverFormWindow.SetContent(c)
	serverFormWindow.Show()
}

func bind(server *Server) {
	serverName := binding.NewString()
	serverName.Set(fmt.Sprintf("服务器名称：%s", "-"))
	playerCount := binding.NewString()
	playerCount.Set(fmt.Sprintf("在线人数：%s", "-"))
	maxDurationInfo := binding.NewString()
	maxDurationInfo.Set(fmt.Sprintf("最长连续在线：%s", "-"))
	remarkInfo := binding.NewString()
	remarkInfo.Set(fmt.Sprintf("备注：%s", server.Remark))

	dataList := binding.BindStringList(&[]string{})

	server.ViewData = &ViewData{
		ServerName:      serverName,
		PlayerCount:     playerCount,
		MaxDurationInfo: maxDurationInfo,
		Remark:          remarkInfo,
		PlayerInfos:     dataList,
	}

	panelContainer := container.NewVBox()

	var scroll *container.Scroll

	overviewContainer := container.NewHBox()
	var toggleBtn *widget.Button
	toggleBtn = widget.NewButton("→", func() {
		if scroll != nil {
			if scroll.Visible() {
				scroll.Hide()
				toggleBtn.SetText("→")
			} else {
				scroll.Show()
				toggleBtn.SetText("↓")
			}
		}
	})
	var editBtn *widget.Button
	editBtn = widget.NewButton("编辑", func() {
		showEditUI(server)
	})
	var removeBtn *widget.Button
	removeBtn = widget.NewButton("-", func() {
		dialog.NewCustomConfirm("提示", "确定", "取消", widget.NewLabel("确定删除吗"), func(b bool) {
			if b {
				serverContainer.RemoveServer(server)
				resetServerConfig()
				err := config.SaveConfig()
				if err != nil {
					dialogutil.ShowInformation("提示", "保存失败", w)
					return
				}
				panelContainer.Hide()
			}
		}, w).Show()
	})
	overviewContainer.Add(toggleBtn)
	overviewContainer.Add(editBtn)
	overviewContainer.Add(removeBtn)
	overviewContainer.Add(widget.NewLabelWithData(serverName))
	overviewContainer.Add(widget.NewLabelWithData(playerCount))
	overviewContainer.Add(widget.NewLabelWithData(maxDurationInfo))
	overviewContainer.Add(widget.NewLabelWithData(remarkInfo))

	panelContainer.Add(overviewContainer)

	list := widget.NewListWithData(dataList, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		s := item.(binding.String)
		o := obj.(*widget.Label)
		o.Bind(s)
		sNew, err := s.Get()
		if err != nil {
			sNew = "-"
		}
		_ = s.Set(sNew)
	})
	scroll = container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(0, 175))
	scroll.Hide()
	detailContainer := container.NewVBox()
	detailContainer.Add(scroll)
	panelContainer.Add(detailContainer)
	serverListPanel.Add(panelContainer)
	serverListPanel.Refresh()
}

func resetServerConfig() {
	serverConfig := make([]map[string]interface{}, 0)
	for _, server := range serverContainer.GetServers() {
		serverConfig = append(serverConfig, map[string]interface{}{
			"ip":       server.Ip,
			"port":     server.Port,
			"interval": server.Interval,
			"remark":   server.Remark,
		})
	}
	viper.Set("servers", serverConfig)
}
