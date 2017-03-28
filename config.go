package gohttp

type Config struct {
    Addr    string
    ReadTimeout int
    WriteTimeout    int
    MaxHeaderBytes  int
    certFile    string
    keyFile string
}

func NewConfig(addr string,read,write,maxheader int) *Config {
    return &Config{
        Addr:addr,
        ReadTimeout:    read,
        WriteTimeout:   write,
        MaxHeaderBytes: maxheader,
    }
}
