package client

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/comoyi/steam-server-monitor/theme"
	a2s "github.com/rumblefrog/go-a2s"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var servers = make([]*Server, 0)
var w fyne.Window
var c *fyne.Container = container.NewVBox()

func Start() {
	log.Debugf("Client start\n")

	versionText := "0.0.1"

	windowTitle := fmt.Sprintf("服务器信息查看器-%s", versionText)

	myApp := app.New()
	myApp.Settings().SetTheme(theme.MyTheme)
	w = myApp.NewWindow(windowTitle)
	w.Resize(fyne.NewSize(400, 600))
	w.SetContent(c)

	addBtn := widget.NewButton("+", func() {
		addWindow := myApp.NewWindow("添加服务器")
		c := container.NewVBox()
		c2 := container.NewAdaptiveGrid(2)
		c3 := container.NewAdaptiveGrid(2)
		c4 := container.NewAdaptiveGrid(2)
		ipLabel := widget.NewLabel("IP")
		ipEntry := widget.NewEntry()
		portLabel := widget.NewLabel("端口")
		portEntry := widget.NewEntry()
		intervalLabel := widget.NewLabel("间隔")
		intervalEntry := widget.NewEntry()
		intervalEntry.Text = "10"
		addBtn := widget.NewButton("添加", func() {
			ip := ipEntry.Text
			if ip == "" {
				dialog.ShowInformation("提示", "请输入IP", addWindow)
				return
			}

			portVal := portEntry.Text
			if portVal == "" {
				dialog.ShowInformation("提示", "请输入端口", addWindow)
				return
			}
			port, err := strconv.ParseInt(portVal, 10, 64)
			if err != nil {
				dialog.ShowInformation("提示", "请输入正确的端口", addWindow)
				return
			}
			if port < 0 {
				dialog.ShowInformation("提示", "请输入正确的端口", addWindow)
				return
			}

			intervalVal := intervalEntry.Text
			if intervalVal == "" {
				dialog.ShowInformation("提示", "请输入间隔", addWindow)
				return
			}
			interval, err := strconv.ParseInt(intervalVal, 10, 64)
			if err != nil {
				dialog.ShowInformation("提示", "请输入正确的间隔", addWindow)
				return
			}
			if interval <= 0 {
				dialog.ShowInformation("提示", "请输入合适的间隔", addWindow)
				return
			}

			server := &Server{
				Name:     "",
				Ip:       ip,
				Port:     port,
				Interval: interval,
				Remark:   "",
				ViewData: nil,
			}
			servers = append(servers, server)
			handleServer(server)

			serverConfig := make([]map[string]interface{}, 0)
			for _, server := range servers {
				serverConfig = append(serverConfig, map[string]interface{}{
					"ip":       server.Ip,
					"port":     server.Port,
					"interval": server.Interval,
				})
			}
			viper.Set("servers", serverConfig)

			addWindow.Close()
		})

		c2.Add(ipLabel)
		c2.Add(ipEntry)
		c3.Add(portLabel)
		c3.Add(portEntry)
		c4.Add(intervalLabel)
		c4.Add(intervalEntry)
		c.Add(c2)
		c.Add(c3)
		c.Add(c4)
		c.Add(addBtn)

		addWindow.Resize(fyne.NewSize(300, 200))
		addWindow.SetContent(c)
		addWindow.Show()
	})
	c.Add(addBtn)

	saveBtn := widget.NewButtonWithIcon("保存", theme.MyTheme.Icon(fynetheme.IconNameDocumentSave), func() {
		log.Debugf("%+v\n", viper.AllSettings())
		config.SaveConfig()
	})
	c.Add(saveBtn)

	for _, s := range config.Conf.Servers {
		server := &Server{
			Ip:       s.Ip,
			Port:     s.Port,
			Interval: s.Interval,
		}
		servers = append(servers, server)
	}

	go func() {
		run()
	}()

	w.ShowAndRun()
}

type Server struct {
	Name     string
	Ip       string
	Port     int64
	Interval int64
	Remark   string
	ViewData *ViewData
}

type ViewData struct {
	ServerName      binding.String
	PlayerCount     binding.String
	MaxDurationInfo binding.String
	PlayerInfos     binding.ExternalStringList
}

type Player struct {
	Duration int64 `json:"duration"`
}

type Info struct {
	ServerName  string    `json:"server_name"`
	PlayerCount int64     `json:"player_count"`
	Players     []*Player `json:"players"`
}

func run() {
	for _, server := range servers {
		handleServer(server)
	}
}

func handleServer(server *Server) {
	bind(server)
	go func(server *Server) {
		var interval int64 = server.Interval
		if interval <= 0 {
			interval = 10
		}
		refresh(server)
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-ticker.C:
				refresh(server)
			}
		}
	}(server)
}

func bind(server *Server) {
	serverName := binding.NewString()
	serverName.Set(fmt.Sprintf("服务器名称：%s", "-"))
	playerCount := binding.NewString()
	playerCount.Set(fmt.Sprintf("在线人数：%s", "-"))
	maxDurationInfo := binding.NewString()
	maxDurationInfo.Set(fmt.Sprintf("最长连续在线：%s", "-"))

	dataList := binding.BindStringList(&[]string{})

	server.ViewData = &ViewData{
		ServerName:      serverName,
		PlayerCount:     playerCount,
		MaxDurationInfo: maxDurationInfo,
		PlayerInfos:     dataList,
	}

	cOverview := container.NewHBox()
	cOverview.Add(widget.NewLabelWithData(serverName))
	cOverview.Add(widget.NewLabelWithData(playerCount))
	cOverview.Add(widget.NewLabelWithData(maxDurationInfo))

	c.Add(cOverview)

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
	scroll := container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(50, 100))
	cDetail := container.NewVBox()
	cDetail.Add(scroll)
	c.Add(cDetail)
}

func refresh(server *Server) {
	info, err := getInfo(server)
	if err != nil {
		return
	}
	refreshUI(server, info)
}

func refreshUI(server *Server, info *Info) {
	infoJson, err := json.Marshal(info)
	if err != nil {
		log.Warnf("json.Marshal failed, err: %v\n", err)
		return
	}
	log.Debugf("infoJson: %s\n", infoJson)

	var maxDuration int64 = 0
	for _, p := range info.Players {
		if p == nil {
			continue
		}
		if p.Duration > maxDuration {
			maxDuration = p.Duration
		}
	}

	server.ViewData.ServerName.Set(fmt.Sprintf("服务器名称：%s", info.ServerName))
	server.ViewData.PlayerCount.Set(fmt.Sprintf("在线人数：%d", info.PlayerCount))
	server.ViewData.MaxDurationInfo.Set(fmt.Sprintf("最长连续在线：%d秒", maxDuration))

	playerInfo := make([]string, 0)
	for i, p := range info.Players {
		if p == nil {
			continue
		}
		playerInfo = append(playerInfo, fmt.Sprintf("玩家%d连续在线%d秒", i+1, p.Duration))
	}

	server.ViewData.PlayerInfos.Set(playerInfo)
}

func getInfo(server *Server) (*Info, error) {
	var err error
	ip := server.Ip
	port := server.Port
	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := a2s.NewClient(address)

	if err != nil {
		log.Warnf("NewClient failed, err: %v\n", err)
		return nil, err
	}

	defer client.Close()

	serverInfo, err := client.QueryInfo()

	if err != nil {
		log.Warnf("QueryInfo failed, err: %v\n", err)
		return nil, err
	}

	serverInfoJson, err := json.Marshal(serverInfo)
	if err != nil {
		log.Warnf("Marshal failed, err: %v\n", err)
		return nil, err
	}
	log.Debugf("serverInfoJson: %s\n", serverInfoJson)

	var serverName = ""
	serverName = serverInfo.Name

	playerInfo, err := client.QueryPlayer()

	if err != nil {
		log.Warnf("QueryPlayer failed, err: %v\n", err)
		return nil, err
	}

	playerInfoJson, err := json.Marshal(playerInfo)
	if err != nil {
		log.Warnf("Marshal failed, err: %v\n", err)
		return nil, err
	}
	log.Debugf("playerInfoJson: %s\n", playerInfoJson)

	var players = make([]*Player, 0)
	for _, p := range playerInfo.Players {
		if p == nil {
			continue
		}
		player := &Player{
			Duration: int64(p.Duration),
		}
		players = append(players, player)
	}

	var playerCount int64 = 0
	playerCount = int64(len(players))
	return &Info{
		ServerName:  serverName,
		PlayerCount: playerCount,
		Players:     players,
	}, nil
}
