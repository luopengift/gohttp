package gohttp

type Config struct {
	Addr              string `json:"addr"`
	ReadTimeout       int    `json:"readtimeout"`
	ReadHeaderTimeout int    `json:"readheadertime"`
	WriteTimeout      int    `json:"writetime"`
	MaxHeaderBytes    int    `json:"maxheaderbytes"`
	CertFile          string `json:"cert"`
	KeyFile           string `json:"key"`
	StaticPath        string `json:"static"`
	LogFormat         string `json:'log_format'`
}

func (cfg *Config) SetLogFormat(format string) {
	cfg.LogFormat = format
}

func (cfg *Config) SetAddress(addr string) {
	cfg.Addr = addr
}

func (cfg *Config) SetTimeout(timeout int) {
	cfg.ReadTimeout = timeout
	cfg.ReadHeaderTimeout = timeout
	cfg.WriteTimeout = timeout
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
	cfg.SetLogFormat("%3d %s %s (%s) %s")
	cfg.SetAddress(":18081")
	cfg.SetTimeout(30)
	cfg.SetMaxHeaderBytes(1 << 20) //DefaultMaxHeaderBytes 1MB
	cfg.SetStaticPath(".")
	return cfg
}
