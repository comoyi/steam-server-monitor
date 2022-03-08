package client

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/comoyi/steam-server-monitor/theme"
	a2s "github.com/rumblefrog/go-a2s"
	"time"
)

var w fyne.Window
var c *fyne.Container = container.NewVBox()

func Start() {
	log.Debugf("Client start\n")

	versionText := "0.0.1"

	windowTitle := fmt.Sprintf("服务器信息查看器-%s", versionText)

	myApp := app.New()
	myApp.Settings().SetTheme(&theme.Theme{})
	w = myApp.NewWindow(windowTitle)
	w.Resize(fyne.NewSize(400, 600))

	go func() {
		refresher()
	}()

	w.ShowAndRun()
}

type Player struct {
	Duration int64 `json:"duration"`
}

type Info struct {
	ServerName  string    `json:"server_name"`
	PlayerCount int64     `json:"player_count"`
	Players     []*Player `json:"players"`
}

func refresher() {
	var interval int64 = 10
	refresh()
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			refresh()
		}
	}
}

func refresh() {
	var err error
	infos, err := getInfos()
	if err != nil {
		log.Warnf("getInfos failed, err: %v\n", err)
		return
	}

	c = container.NewVBox()

	for _, info := range infos {
		if info == nil {
			continue
		}
		handleOne(info)
	}

}

func handleOne(info *Info) {
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

	serverName := binding.NewString()
	serverName.Set(fmt.Sprintf("服务器名称：%s", info.ServerName))
	playerCount := binding.NewString()
	playerCount.Set(fmt.Sprintf("在线人数：%d", info.PlayerCount))
	maxDurationInfo := binding.NewString()
	maxDurationInfo.Set(fmt.Sprintf("最长连续在线：%d", maxDuration))

	//c := container.NewVBox()
	c.Resize(fyne.NewSize(400, 600))
	c.Add(widget.NewLabelWithData(serverName))
	c.Add(widget.NewLabelWithData(playerCount))
	c.Add(widget.NewLabelWithData(maxDurationInfo))
	w.SetContent(c)
}

func getInfos() ([]*Info, error) {
	infos := make([]*Info, 0)
	for _, s := range config.Conf.Servers {
		info, err := getInfo(s)
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func getInfo(server config.Server) (*Info, error) {
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
