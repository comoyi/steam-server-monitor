package config

var Conf Config

type Config struct {
	Debug   bool     `toml:"debug"`
	Servers []Server `toml:"servers"`
}

type Server struct {
	Ip   string `toml:"ip"`
	Port int64  `toml:"port"`
}
