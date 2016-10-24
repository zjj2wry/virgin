package virgin

import "net/http"

type Virgin struct {
	Route  *Route
	Server *http.Server
}

func NewVirgin() *Virgin {
	route := NewRoute()
	Virgin := &Virgin{Route: route, Server: &http.Server{}}

	return Virgin
}

func (virgin *Virgin) AddRoute(method, pattern string, handleFunc HandleFunc) {
	virgin.Route.Add(method, pattern, handleFunc)
}

func (virgin *Virgin) Listen(port string) error {
	virgin.Server.Handler = virgin.Route
	virgin.Server.Addr = port
	
	return virgin.Server.ListenAndServe()
}
