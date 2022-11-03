package client

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/comoyi/steam-server-monitor/util/timeutil"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rumblefrog/go-a2s"
	"sync"
	"time"
)

var serverContainer = NewServerContainer()

type ServerContainer struct {
	Servers []*Server
	mu      sync.Mutex
}

func NewServerContainer() *ServerContainer {
	servers := make([]*Server, 0)
	return &ServerContainer{
		Servers: servers,
		mu:      sync.Mutex{},
	}
}

func (sc *ServerContainer) GetServers() []*Server {
	return sc.Servers
}

func (sc *ServerContainer) AddServer(server *Server) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.Servers = append(sc.Servers, server)
}

func (sc *ServerContainer) RemoveServer(server *Server) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	for i, s := range sc.Servers {
		if s == server {
			sc.Servers = append(sc.Servers[:i], sc.Servers[i+1:]...)
			break
		}
	}
}

type Server struct {
	DisplayName    string
	Name           string
	Ip             string
	Port           int64
	Interval       int64
	IntervalTicker *time.Ticker
	Remark         string
	Info           *Info
	ViewData       *ViewData
}

func NewServer(displayName string, ip string, port int64, interval int64, remark string) *Server {
	if interval <= 0 {
		interval = 10
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	return &Server{
		DisplayName:    displayName,
		Ip:             ip,
		Port:           port,
		Interval:       interval,
		IntervalTicker: ticker,
		Remark:         remark,
	}
}

func (s *Server) Start() {
	s.AsyncRefresh()
}

func (s *Server) AsyncRefresh() {
	go func(server *Server) {
		refresh(server)
		for {
			select {
			case <-server.IntervalTicker.C:
				refresh(server)
			}
		}
	}(s)
}

func (s *Server) UpdateInterval(interval int64) {
	s.Interval = interval
	if s.IntervalTicker != nil {
		s.IntervalTicker.Reset(time.Duration(interval) * time.Second)
	}
}

func (s *Server) getInfo() (*Info, error) {
	return getInfo(s)
}

type ViewData struct {
	ServerName      binding.String
	PlayerCount     binding.String
	MaxDurationInfo binding.String
	Remark          binding.String
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

func refresh(server *Server) {
	info, err := server.getInfo()
	if err != nil {
		return
	}
	server.Info = info
	refreshUI(server)
}

func refreshUI(server *Server) {
	if server == nil {
		log.Warnf("refreshUI server is nil\n")
		return
	}
	info := server.Info
	infoJson, err := json.Marshal(info)
	if err != nil {
		log.Warnf("json.Marshal failed, err: %v\n", err)
		return
	}
	log.Debugf("infoJson: %s\n", infoJson)

	if server.DisplayName != "" {
		server.ViewData.ServerName.Set(fmt.Sprintf("服务器名称：%s", server.DisplayName))
	}
	server.ViewData.Remark.Set(fmt.Sprintf("备注：%s", server.Remark))

	if info != nil {
		var maxDuration int64 = 0
		for _, p := range info.Players {
			if p == nil {
				continue
			}
			if p.Duration > maxDuration {
				maxDuration = p.Duration
			}
		}
		maxDurationFormatted := "-"
		if info.PlayerCount > 0 {
			maxDurationFormatted = timeutil.FormatDuration(maxDuration)
		}

		serverNameFixed := ""
		if server.DisplayName != "" {
			serverNameFixed = server.DisplayName
		} else {
			serverNameFixed = bluemonday.StrictPolicy().Sanitize(info.ServerName)
		}
		server.ViewData.ServerName.Set(fmt.Sprintf("服务器名称：%s", serverNameFixed))
		server.ViewData.PlayerCount.Set(fmt.Sprintf("在线人数：%d", info.PlayerCount))
		server.ViewData.MaxDurationInfo.Set(fmt.Sprintf("最长连续在线：%s", maxDurationFormatted))

		playerInfoList := make([]string, 0)
		for i, p := range info.Players {
			if p == nil {
				continue
			}
			playerInfoList = append(playerInfoList, fmt.Sprintf("玩家%d连续在线%s", i+1, timeutil.FormatDuration(p.Duration)))
		}

		server.ViewData.PlayerInfos.Set(playerInfoList)
	}
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
