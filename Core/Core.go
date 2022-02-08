package Core

import (
	"log"

	"gopkg.in/ini.v1"
)

// var Config = struct {
// 	HTTP struct {
// 		Port           string
// 		Token          string
// 		ConsoleDisable bool
// 	}
// 	Dns struct {
// 		Domain string
// 		Xip    string
// 		Dnslog string
// 	}
// }{}
type ConfigVar struct {
	HTTP HttpConfig `ini:"HTTP"`
	Dns  DnsConfig  `ini:"DNS"`
}

type HttpConfig struct {
	Port           string `ini:"Port"`
	Token          string `ini:"Token"`
	ConsoleDisable bool   `ini:"ConsoleDisable"`
}
type DnsConfig struct {
	Domain string `ini:"Domain"`
	Xip    string `ini:"Xip"`
	Dnslog string `ini:"Dnslog"`
}

var Config ConfigVar
var cfg, err = ini.Load("./config.ini")

func (c *ConfigVar) Init() {
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg.MapTo(&Config)
}

func (c *ConfigVar) GetConfig() ConfigVar {
	return Config
}

func (c *ConfigVar) GetCfg() *ini.File {
	return cfg
}

func (c *ConfigVar) SaveCfg() {
	cfg.SaveTo("./config.ini")
}
