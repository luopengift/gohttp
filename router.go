package gohttp

import (
	"reflect"
	"regexp"
)

type muxEntry reflect.Type

type route struct {
	path  string
	regx  *regexp.Regexp
	entry muxEntry
}

func newroute(path string, handler Handler) *route {
	rv := reflect.ValueOf(handler)
	rt := reflect.Indirect(rv).Type()
	return &route{path: path, regx: regexp.MustCompile(path), entry: rt}
}

type RouterList []*route

func InitRouterList() *RouterList {
	return new(RouterList)
}

func (self *RouterList) Route(path string, handler Handler) {
	route := newroute(path, handler)
	*self = append(*self, route)
}

func (self *RouterList) Find(path string) (muxEntry, map[string]string) {
	for _, route := range *self {
		if match := route.regx.FindStringSubmatch(path); match != nil {
			kv := make(map[string]string)
			for key, value := range route.regx.SubexpNames() {
				kv[value] = match[key]
			}
			delete(kv, "")
			return route.entry, kv
		}
	}
	return nil, nil
}

func (self *RouterList) String() string {
	str := "\nRouter Map:\n"
	for _, route := range *self {
		str += route.path + " => " + route.entry.String() + "\n"
	}
	return str
}
