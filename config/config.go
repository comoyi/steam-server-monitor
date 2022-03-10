package config

import (
	"fmt"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var Conf Config

type Config struct {
	Debug   bool     `toml:"debug"`
	Servers []Server `toml:"servers"`
}

type Server struct {
	Ip       string `toml:"ip"`
	Port     int64  `toml:"port"`
	Interval int64  `toml:"interval"`
	Remark   string `toml:"remark"`
}

func LoadConfig() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("%s%s%s", "$HOME", string(os.PathSeparator), ".steam-server-monitor"))
	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("Read config failed, err: %v\n", err)
		return
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Errorf("Unmarshal config failed, err: %v\n", err)
		return
	}
	log.Debugf("config: %+v\n", Conf)
}

var saveMutex = &sync.Mutex{}

func SaveConfig() error {
	saveMutex.Lock()
	defer saveMutex.Unlock()
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("Get os.UserHomeDir failed, err: %v\n", err)
		return err
	}
	log.Debugf("userHomeDir: %s\n", userHomeDir)

	configPath := fmt.Sprintf("%s%s%s", userHomeDir, string(os.PathSeparator), ".steam-server-monitor")
	configFile := fmt.Sprintf("%s%s%s", configPath, string(os.PathSeparator), "config.toml")
	log.Debugf("configFile: %s\n", configFile)

	exist, err := isPathExist(configPath)
	if err != nil {
		log.Warnf("Check isPathExist failed, err: %v\n", err)
		return err
	}
	if !exist {
		err = os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			log.Warnf("Get os.Mkdir failed, err: %v\n", err)
			return err
		}
	}

	err = viper.WriteConfigAs(configFile)
	if err != nil {
		log.Errorf("SafeWriteConfigAs failed, err: %v\n", err)
		return err
	}
	return nil
}

func isPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}
