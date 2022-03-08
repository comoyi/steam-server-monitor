package config

type Config struct {
	Ip   string `yaml:"ip"`
	Port int64  `yaml:"port"`
}

var Conf Config
