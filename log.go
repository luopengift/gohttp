package gohttp

import (
	"github.com/luopengift/log"
	"os"
)

// InitLog inits gohttp loghandler
func InitLog() *log.Log {
	fileHandler := log.NewFile("/tmp/access_%Y%M%D.log")
	gohttpLog := log.NewLog("gohttp", os.Stderr, fileHandler)
	gohttpLog.SetFormatter(log.NewTextFormat("TIME [LEVEL] MESSAGE", 0))
	return gohttpLog
}
