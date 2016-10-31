package virgin

import "net/http"

type Virgin struct {
	Router *Router
	Server *http.Server
}

func NewVirgin() *Virgin {
	Router := NewRouter()
	Virgin := &Virgin{Router: Router, Server: &http.Server{}}

	return Virgin
}

func (virgin *Virgin) AddRoute(method, pattern string, handlerFunc HandlerFunc) {
	virgin.Router.Add(method, pattern, handlerFunc)
}

func (virgin *Virgin) Listen(port string) error {
	virgin.Server.Handler = virgin.Router
	virgin.Server.Addr = port

	return virgin.Server.ListenAndServe()
}
