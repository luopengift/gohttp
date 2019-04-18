package gohttp

import (
	"os"

	"github.com/luopengift/log"
)

// Logger interface
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// InitLog inits gohttp loghandler
func InitLog() *log.Log {
	//fileHandler := log.NewFile("/tmp/access_%Y%M%D.log")
	gohttpLog := log.NewLog("gohttp", os.Stderr) //, fileHandler)
	gohttpLog.SetFormatter(log.NewTextFormat("TIME [LEVEL] MESSAGE", log.ModeColor))
	return gohttpLog
}
