package api

import (
	"encoding/json"
	"fmt"
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/comoyi/steam-server-monitor/config"
	"github.com/comoyi/steam-server-monitor/log"
	"net/http"
	"strconv"
)

func Start() {

	http.HandleFunc("/player-count", playerCount)
	http.HandleFunc("/info", info)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Conf.ApiPort), nil)
	if err != nil {
		fmt.Printf("server start failed err: %v\n", err)
		log.Errorf("server start failed err: %v\n", err)
		return
	}
}

func playerCount(writer http.ResponseWriter, request *http.Request) {
	var err error
	host := request.URL.Query().Get("host")
	portRaw, _ := strconv.Atoi(request.URL.Query().Get("port"))
	var port int64 = int64(portRaw)
	serverContainer := client.GetServerContainer()
	servers := serverContainer.GetServers()
	var server *client.Server
	for _, s := range servers {
		if s.Ip == host && s.Port == port {
			server = s
			break
		}
	}

	var num int64 = -1
	if server != nil {
		if server.Info != nil {
			num = server.Info.PlayerCount
		}
	}

	bytes := []byte(fmt.Sprintf("%v", num))
	writer.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	_, err = writer.Write(bytes)
	if err != nil {
		log.Debugf("write failed, err: %s\n", err)
		return
	}
}

func info(writer http.ResponseWriter, request *http.Request) {

	var err error
	host := request.URL.Query().Get("host")
	portRaw, _ := strconv.Atoi(request.URL.Query().Get("port"))
	var port int64 = int64(portRaw)
	serverContainer := client.GetServerContainer()
	servers := serverContainer.GetServers()
	var server *client.Server
	for _, s := range servers {
		if s.Ip == host && s.Port == port {
			server = s
			break
		}
	}

	bytes := []byte("")
	if server != nil {
		bytes, err = json.Marshal(server.Info)
	}
	if err != nil {
		log.Debugf("json.Marshal failed, err: %s\n", err)
		return
	}

	j := string(bytes)
	log.Debugf("json: %s\n", j)
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = writer.Write(bytes)
	if err != nil {
		log.Debugf("write failed, err: %s\n", err)
		return
	}
}
