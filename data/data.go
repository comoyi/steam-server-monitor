package data

import (
	"encoding/json"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func New() *Data {
	data := &Data{}
	return data
}

type Data struct {
	Counter   int
	ChCounter chan struct{}
	Servers   []*Server
}

type Server struct {
	DisplayName string
	Ip          string
	Port        int64
	Interval    int64
	Remark      string
}

func (d *Data) Init() error {
	d.ChCounter = make(chan struct{})
	d.Servers = make([]*Server, 0)

	conf, err := LoadConfig()
	if err != nil {
		log.Errorf("load server config failed, err: %v", err)
		return err
	}

	serverConfigJsonBytes, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	serverConfigJson := string(serverConfigJsonBytes)
	log.Debugf("load server config success, conf: %v", serverConfigJson)

	for _, v := range conf.Servers {
		server := &Server{
			DisplayName: v.DisplayName,
			Ip:          v.Ip,
			Port:        v.Port,
			Interval:    v.Interval,
			Remark:      v.Remark,
		}
		d.Servers = append(d.Servers, server)
	}

	return nil
}

type Config struct {
	Servers []*ServerConfig
}

type ServerConfig struct {
	DisplayName string
	Ip          string
	Port        int64
	Interval    int64
	Remark      string
}

func LoadConfig() (*Config, error) {
	var conf *Config

	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fixedExePath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(fixedExePath)

	v := viper.New()
	v.SetConfigName("server")
	v.SetConfigType("toml")
	v.AddConfigPath(exeDir)
	v.AddConfigPath(filepath.Join(exeDir, "config"))

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	usedConfig := v.ConfigFileUsed()
	log.Debugf("used server config: %s", usedConfig)

	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
