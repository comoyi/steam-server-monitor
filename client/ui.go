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
	"runtime"
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

	myApp = app.NewWithID("com.comoyi.steamservermonitor")
	myApp.Settings().SetTheme(theme.CustomTheme)
	w = myApp.NewWindow(windowTitle)
	w.SetMaster()
	w.Resize(fyne.NewSize(400, 600))
	c = container.NewVBox()
	w.SetContent(c)

	if runtime.GOOS == "android" {
		hc := container.NewCenter()
		hc.Add(widget.NewLabel(windowTitle))
		c.Add(hc)
	}

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

		authorInfo := widget.NewLabel("Copyright © 2022-2023 清新池塘")
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

var serverFormWindow fyne.Window

func showServerFormUI(isEdit bool, server *Server) {
	title := "添加服务器"
	if isEdit {
		if server == nil {
			return
		}
		title = "编辑服务器"
	}

	if serverFormWindow != nil {
		// prevent error exit on android
		if runtime.GOOS != "android" {
			serverFormWindow.Close()
		}
	}
	serverFormWindow = myApp.NewWindow(title)

	c := container.NewVBox()
	c1 := container.NewAdaptiveGrid(2)
	c2 := container.NewAdaptiveGrid(2)
	c3 := container.NewAdaptiveGrid(2)
	c4 := container.NewAdaptiveGrid(2)
	c5 := container.NewAdaptiveGrid(2)

	displayNameLabel := widget.NewLabel("显示名称")
	var displayNameEntry *widget.Entry
	displayNameEntry = widget.NewEntry()
	displayNameEntry.SetPlaceHolder("默认为服务器名称")
	if isEdit {
		displayNameEntry.SetText(server.DisplayName)
	}

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
	displayName := displayNameEntry.Text
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
			server.DisplayName = displayName
			server.Ip = ip
			server.Port = port
			server.UpdateInterval(interval)
			server.Remark = remark
			refreshUI(server)
		} else {
			newServer := NewServer(displayName, ip, port, interval, remark)
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

		serverFormWindow.Close()
	})
	submitBtn.SetIcon(theme2.DocumentSaveIcon())

	var removeBtn *widget.Button
	removeBtn = widget.NewButtonWithIcon("", theme2.DeleteIcon(), func() {
		dialog.NewCustomConfirm("提示", "确定", "取消", widget.NewLabel(fmt.Sprintf("确定删除吗\n%s", displayName)), func(b bool) {
			if b {
				serverContainer.RemoveServer(server)
				resetServerConfig()
				err := config.SaveConfig()
				if err != nil {
					dialogutil.ShowInformation("提示", "保存失败", serverFormWindow)
					return
				}

				// remove UI container
				serverListPanel.Remove(server.Container)

				serverFormWindow.Close()
			}
		}, serverFormWindow).Show()
	})
	if !isEdit {
		removeBtn.Disable()
	}

	c1.Add(displayNameLabel)
	c1.Add(displayNameEntry)
	c2.Add(ipLabel)
	c2.Add(ipEntry)
	c3.Add(portBox)
	c3.Add(portEntry)
	c4.Add(intervalLabel)
	c4.Add(intervalEntry)
	c5.Add(remarkLabel)
	c5.Add(remarkEntry)
	c.Add(c1)
	c.Add(c2)
	c.Add(c3)
	c.Add(c4)
	c.Add(c5)
	cop1 := container.NewGridWithColumns(2)
	cop2 := container.NewVBox()
	cop3 := container.NewVBox()
	cop1.Add(cop2)
	cop2.Add(removeBtn)
	cop1.Add(cop3)
	cop3.Add(submitBtn)
	c.Add(cop1)

	serverFormWindow.SetContent(c)
	serverFormWindow.Show()
}

func bind(server *Server) {
	serverName := binding.NewString()
	displayName := "-"
	if server.DisplayName != "" {
		displayName = server.DisplayName
	}
	serverName.Set(fmt.Sprintf("服务器：%s", displayName))
	playerCount := binding.NewString()
	playerCount.Set(fmt.Sprintf("在线人数：%s", "-"))
	maxDurationInfo := binding.NewString()
	maxDurationInfo.Set(fmt.Sprintf("最长在线：%s", "-"))
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
	server.Container = panelContainer

	var detailContainer *fyne.Container

	var toggleBtn *widget.Button
	toggleBtn = widget.NewButton("→", func() {
		if detailContainer != nil {
			if detailContainer.Visible() {
				detailContainer.Hide()
				toggleBtn.SetText("→")
			} else {
				detailContainer.Show()
				toggleBtn.SetText("↓")
			}
		}
	})
	var editBtn *widget.Button
	editBtn = widget.NewButton("", func() {
		showEditUI(server)
	})
	editBtn.SetIcon(theme2.DocumentCreateIcon())

	overviewContainer := container.NewHBox()
	b1 := container.NewVBox()
	overviewContainer.Add(b1)
	b2 := container.NewHBox()
	b3 := container.NewHBox()
	detailContainer = container.NewHBox()
	detailContainer.Hide()
	b7 := container.NewHBox()
	b1.Add(b2)
	b1.Add(b3)
	b1.Add(detailContainer)
	b1.Add(b7)
	b4 := container.NewVBox()
	b5 := container.NewVBox()
	b3.Add(toggleBtn)
	b3.Add(b4)
	b3.Add(b5)
	b2.Add(editBtn)
	b2.Add(widget.NewLabelWithData(serverName))
	b4.Add(widget.NewLabelWithData(playerCount))
	b5.Add(widget.NewLabelWithData(maxDurationInfo))

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

	var scroll *container.Scroll
	scroll = container.NewVScroll(list)
	detailListContainer := container.NewGridWrap(fyne.NewSize(320, 190))
	detailListContainer.Add(scroll)

	detailContainer.Add(container.NewGridWrap(fyne.NewSize(40, 40)))
	detailContainer.Add(detailListContainer)
	if server.Remark != "" {
		b7.Add(container.NewGridWrap(fyne.NewSize(40, 40)))
		b7.Add(widget.NewLabelWithData(remarkInfo))
	}

	panelContainer.Add(overviewContainer)
	serverListPanel.Add(panelContainer)
	serverListPanel.Refresh()
}

func resetServerConfig() {
	serverConfig := make([]map[string]interface{}, 0)
	for _, server := range serverContainer.GetServers() {
		serverConfig = append(serverConfig, map[string]interface{}{
			"display_name": server.DisplayName,
			"ip":           server.Ip,
			"port":         server.Port,
			"interval":     server.Interval,
			"remark":       server.Remark,
		})
	}
	viper.Set("servers", serverConfig)
}
