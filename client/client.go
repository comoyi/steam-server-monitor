package client

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/comoyi/steam-server-monitor/log"
	a2s "github.com/rumblefrog/go-a2s"
	"time"
)

func Start() {
	log.Debugf("Client start\n")

	versionText := "0.0.1"

	windowTitle := fmt.Sprintf("服务信息查看器-%s", versionText)

	myApp := app.New()
	w := myApp.NewWindow(windowTitle)
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
	ServerName   string    `json:"server_name"`
	PlayersCount int64     `json:"players_count"`
	Players      []*Player `json:"players"`
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
		ServerName:   serverName,
		PlayersCount: playerCount,
		Players:      players,
	}, nil
}
