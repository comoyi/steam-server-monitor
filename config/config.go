package config

import (
	"fyne.io/fyne/v2/app"
	"github.com/comoyi/steam-server-monitor/log"
	"github.com/comoyi/steam-server-monitor/util/fsutil"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var Conf Config

type Config struct {
	LogLevel string    `toml:"log_level" mapstructure:"log_level"`
	Servers  []*Server `toml:"servers" mapstructure:"servers"`
}

type Server struct {
	DisplayName string `toml:"display_name" mapstructure:"display_name"`
	Ip          string `toml:"ip" mapstructure:"ip"`
	Port        int64  `toml:"port" mapstructure:"port"`
	Interval    int64  `toml:"interval" mapstructure:"interval"`
	Remark      string `toml:"remark" mapstructure:"remark"`
}

func initDefaultConfig() {
	viper.SetDefault("log_level", log.Off)
}

func LoadConfig() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	configDirPath, err := getConfigDirPath()
	if err != nil {
		log.Warnf("Get configDirPath failed, err: %v\n", err)
		return
	}
	viper.AddConfigPath(configDirPath)

	initDefaultConfig()

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

	err := viper.WriteConfig()
	if err == nil {
		return nil
	}

	configDirPath, err := getConfigDirPath()
	if err != nil {
		log.Warnf("Get configDirPath failed, err: %v\n", err)
		return err
	}
	configFile := filepath.Join(configDirPath, "config.toml")
	log.Debugf("configFile: %s\n", configFile)

	exist, err := fsutil.Exists(configDirPath)
	if err != nil {
		log.Warnf("Check isPathExist failed, err: %v\n", err)
		return err
	}
	if !exist {
		err = os.MkdirAll(configDirPath, os.ModePerm)
		if err != nil {
			log.Warnf("Get os.MkdirAll failed, err: %v\n", err)
			return err
		}
	}

	err = viper.WriteConfigAs(configFile)
	if err != nil {
		log.Errorf("WriteConfigAs failed, err: %v\n", err)
		return err
	}
	return nil
}

func getConfigDirPath() (string, error) {
	configRootPath, err := getConfigRootPath()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configRootPath, ".steam-server-monitor")
	return configPath, nil
}

func getConfigRootPath() (string, error) {
	var err error
	configRootPath := ""
	if runtime.GOOS == "android" {
		configRootPath = app.NewWithID("com.comoyi.steamservermonitor").Storage().RootURI().Path()
	} else {
		configRootPath, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	return configRootPath, nil
}
