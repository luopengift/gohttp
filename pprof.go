package gohttp

import (
	"net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
)

var (
	// StartTrace 运行trace
	StartTrace = func(ctx *Context) {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
	}

	// StopTrace 停止trace
	StopTrace = func(ctx *Context) {
		trace.Stop()
	}

	// StartGC 手动触发GC
	StartGC = func(ctx *Context) {
		runtime.GC()
	}

	// StopGC stop gc
	StopGC = func(ctx *Context) {
		debug.SetGCPercent(-1)
	}

	// Index wapper pprof.Index, default path is /debug/pprof/
	Index = pprof.Index

	// Cmdline wapper pprof.Cmdline, default path is /debug/pprof/cmdline
	Cmdline = pprof.Cmdline

	// Profile wapper pprof.Profile, default path is /debug/pprof/profile
	Profile = pprof.Profile

	// Symbol wapper pprof.Symbol, default path is /debug/pprof/symbol
	Symbol = pprof.Symbol

	// Trace wapper pprof.Trace, default path is /debug/pprof/trace
	Trace = pprof.Trace
)
