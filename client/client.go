package client

import (
	"encoding/json"
	"fmt"
	"github.com/comoyi/steam-server-monitor/data"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/rumblefrog/go-a2s"
	"time"
)

func New() *Client {
	client := &Client{}
	return client
}

type Client struct {
	Data *data.Data
}

func (c *Client) Run() {
	go func() {
		for {
			time.Sleep(time.Second)
			c.Data.Counter++
			c.Data.ChCounter <- struct{}{}
		}
	}()

	go func() {
		for _, v := range c.Data.Servers {
			go func(v *data.Server) {
				for {
					select {
					case <-time.After(time.Second * time.Duration(v.Interval)):
						info, err := QueryInfo(v)
						if err != nil {
							log.Infof("Query info error: %v", err)
							return
						}
						infoJsonBytes, err := json.Marshal(info)
						if err != nil {
							return
						}
						infoJson := string(infoJsonBytes)
						log.Debugf("info: %v", string(infoJson))
					}
				}
			}(v)
		}
	}()
}

type Info struct {
	ServerName  string
	PlayerCount int64
	Players     []*Player
}

type Player struct {
	Name     string
	Duration int64
}

func QueryInfo(server *data.Server) (*Info, error) {
	ip := server.Ip
	port := server.Port
	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := a2s.NewClient(address)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	serverInfo, err := client.QueryInfo()
	if err != nil {
		return nil, err
	}
	serverInfoJsonBytes, err := json.Marshal(serverInfo)
	if err != nil {
		return nil, err
	}
	serverInfoJson := string(serverInfoJsonBytes)
	log.Debugf("serverInfoJson: %v", serverInfoJson)
	serverName := serverInfo.Name

	playerInfo, err := client.QueryPlayer()

	if err != nil {
		return nil, err
	}

	playerInfoJsonBytes, err := json.Marshal(playerInfo)
	if err != nil {
		return nil, err
	}
	playerInfoJson := string(playerInfoJsonBytes)
	log.Debugf("playerInfoJson: %v", playerInfoJson)

	var players = make([]*Player, 0)
	for _, p := range playerInfo.Players {
		if p == nil {
			continue
		}
		player := &Player{
			Name:     p.Name,
			Duration: int64(p.Duration),
		}
		players = append(players, player)
	}

	var playerCount int64 = 0
	playerCount = int64(len(players))

	info := &Info{
		ServerName:  serverName,
		PlayerCount: playerCount,
		Players:     players,
	}
	return info, nil
}
