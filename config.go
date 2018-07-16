package gohttp

// Config config
type Config struct {
	Addr              string `json:"addr"`
	ReadTimeout       int    `json:"readtimeout"`
	ReadHeaderTimeout int    `json:"readheadertime"`
	WriteTimeout      int    `json:"writetime"`
	MaxHeaderBytes    int    `json:"maxheaderbytes"`
	CertFile          string `json:"cert"`
	KeyFile           string `json:"key"`
	WebPath           string `json:"web_path"`
	StaticPath        string `json:"static_path"`
	LogFormat         string `json:"log_format"`
}

// SetLogFormat set log format
func (cfg *Config) SetLogFormat(format string) {
	cfg.LogFormat = format
}

// SetAddress set address
func (cfg *Config) SetAddress(addr string) {
	cfg.Addr = addr
}

// SetTimeout set time out
func (cfg *Config) SetTimeout(timeout int) {
	cfg.ReadTimeout = timeout
	cfg.ReadHeaderTimeout = timeout
	cfg.WriteTimeout = timeout
}

// SetMaxHeaderBytes set max header bytes
func (cfg *Config) SetMaxHeaderBytes(max int) {
	cfg.MaxHeaderBytes = max
}

// SetTLS set tls
func (cfg *Config) SetTLS(cert, key string) {
	cfg.CertFile = cert
	cfg.KeyFile = key
}

// SetWebPath set web path
func (cfg *Config) SetWebPath(web string) {
	cfg.WebPath = web
}

// SetStaticPath set static path
func (cfg *Config) SetStaticPath(static string) {
	cfg.StaticPath = static
}

// InitConfig init config
func InitConfig() *Config {
	cfg := new(Config)
	cfg.SetLogFormat("%3d %s %s (%s) %s")
	cfg.SetAddress(":18081")
	cfg.SetTimeout(30)
	cfg.SetMaxHeaderBytes(1 << 20) //DefaultMaxHeaderBytes 1MB
	cfg.SetStaticPath(".")
	return cfg
}
