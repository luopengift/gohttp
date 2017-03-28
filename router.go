package gohttp

import (
	"reflect"
	"regexp"
)

type muxEntry reflect.Type
var Router = map[*regexp.Regexp]muxEntry{}

func RouterRegister(path string, handler Handler) {
	rv := reflect.ValueOf(handler)
	rt := reflect.Indirect(rv).Type()
	Router[regexp.MustCompile(path)] = rt
}
