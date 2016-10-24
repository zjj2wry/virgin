package virgin

import (
	"log"
	"net/http"
)

type HandleFunc func(*Context)

type HandleMap map[string]HandleFunc

var handleMap HandleMap

type Route struct {
	Tree map[string]*HandleMap
}

func NewRoute() (route *Route) {
	route = &Route{
		make(map[string]*HandleMap),
	}

	return route
}

func (r *Route) Add(method, pattern string, handleFunc HandleFunc) {
	if handleMap == nil {
		handleMap = make(HandleMap)
	}

	handleMap[pattern] = handleFunc

	r.Tree[method] = &handleMap
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, re *http.Request) {
	method := re.Method
	uri := re.URL.Path
	log.Printf("%s %s", method, uri)
	handleMap := r.Tree["GET"]
	handfunc := (*handleMap)["/a"]
	ctx := &Context{
		re,
		rw,
	}
	handfunc(ctx)
}
