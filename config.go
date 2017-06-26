package gohttp

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"time"
)

type Config struct {
	Addr              string `json:"addr"`
	ReadTimeout       int    `json:"readtimeout"`
	ReadHeaderTimeout int    `json:"readheadertime"`
	WriteTimeout      int    `json:"writetime"`
	MaxHeaderBytes    int    `json:"maxheaderbytes"`
	CertFile          string ``
	KeyFile           string
	StaticPath        string
}

func (cfg *Config) SetAddress(addr string) {
	cfg.Addr = addr
}

func (cfg *Config) SetTimeout(timeout int) {
	cfg.ReadTimeout = int(time.Duration(timeout) * time.Second)
	cfg.ReadHeaderTimeout = int(time.Duration(timeout) * time.Second)
	cfg.WriteTimeout = int(time.Duration(timeout) * time.Second)
}

func (cfg *Config) SetMaxHeaderBytes(max int) {
	cfg.MaxHeaderBytes = max
}

func (cfg *Config) SetSSL(cert, key string) {
	cfg.CertFile = cert
	cfg.KeyFile = key
}

func (cfg *Config) SetStaticPath(path string) {
	cfg.StaticPath = path
}

func InitConfig() *Config {
	cfg := new(Config)
	cfg.SetAddress(":18081")
	cfg.SetTimeout(30)
	cfg.SetMaxHeaderBytes(1 << 20) //DefaultMaxHeaderBytes 1MB
	return cfg
}

func NewConfig() *Config {
	return InitConfig()
}

func ReadConfig(name string) *Config {
	if name == "" {
		name = "./config.json"
	}
	cfg := new(Config)
	conf := file.NewConfig(name)
	err := conf.Parse(cfg)
	if err != nil {
		panic(err)
	}
	logger.Info("%+v", conf.String())
	return cfg
}
