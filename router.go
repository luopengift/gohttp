package gohttp

import (
	"reflect"
	"regexp"
)

type muxEntry reflect.Type

var RouterMap = map[*regexp.Regexp]muxEntry{}

func RouterRegister(path string, handler Handler) {
	rv := reflect.ValueOf(handler)
	rt := reflect.Indirect(rv).Type()
	RouterMap[regexp.MustCompile(path)] = rt
}
