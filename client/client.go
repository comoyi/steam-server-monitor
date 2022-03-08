package client

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/comoyi/steam-server-monitor/log"
	a2s "github.com/rumblefrog/go-a2s"
	"time"
)

var w fyne.Window

func Start() {
	log.Debugf("Client start\n")

	versionText := "0.0.1"

	windowTitle := fmt.Sprintf("服务器信息查看器-%s", versionText)

	myApp := app.New()
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
	info, err := getInfo()
	if err != nil {
		log.Warnf("getInfo failed, err: %v\n", err)
		return
	}

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
	playerCount.Set(fmt.Sprintf("最长连续在线：%d", maxDuration))

	w.SetContent(container.NewVBox(
		widget.NewLabelWithData(serverName),
		widget.NewLabelWithData(playerCount),
		widget.NewLabelWithData(maxDurationInfo),
	))
}

func getInfo() (*Info, error) {
	var err error
	ip := "127.0.0.1"
	port := 2457
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
