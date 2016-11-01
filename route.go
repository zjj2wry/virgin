package virgin

import (
	"log"
	"net/http"
)

// "log"

type HandlerFunc func(*Context)

type (
	Router struct {
		tree map[string]*node
	}
	node struct {
		label       byte
		prefix      string
		child       []*node
		paramname   string
		handlerFunc HandlerFunc
	}
)

func NewRouter() *Router {
	route := &Router{
		tree: make(map[string]*node),
	}
	return route
}

func (r *Router) Add(method, path string, h HandlerFunc) {
	if path == "" {
		panic("virgin: path cannot be empty")
	}
	if path[0] != '/' {
		path = "/" + path
	}

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			if path[i-1] == '/' {
				r.insert(method, path[:i], nil, "")
				r.insert(method, path[i:], h, path[i+1:])
				return
			}
		} else if path[i] == '*' {
			if path[i-1] == '/' {
				r.insert(method, path[:i], nil, "")
				r.insert(method, path[i:], h, "")
				return
			}
		}
	}

	r.insert(method, path, h, "")
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func (n *node) Add(method, path string, h HandlerFunc, paramname string) {
	search := path

	for {
		sl := len(search)
		pl := len(n.prefix)
		l := 0

		max := pl
		if sl < max {
			max = sl
		}

		for ; l < max && search[l] == n.prefix[l]; l++ {
		}

		if l == 0 {
			n.label = search[0]
			n.prefix = search
			if h != nil {
				n.handlerFunc = h
				n.paramname = paramname
			}
		} else if l < pl {
			n1 := &node{
				n.label,
				n.prefix,
				n.child,
				n.paramname,
				n.handlerFunc,
			}

			n.label = n.prefix[0]
			n.prefix = n.prefix[:l]

			n.child = append(n.child, n1)

			if l == sl {
				n.handlerFunc = h
				n.paramname = paramname
			} else {
				prefix := search[l:]
				n2 := &node{
					prefix[0],
					prefix,
					nil,
					paramname,
					h,
				}
				n.child = append(n.child, n2)
			}
		} else if l < sl {
			search = search[l:]
			c := n.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				n = c
				continue
			}
			n1 := &node{
				search[0],
				search,
				nil,
				paramname,
				h,
			}
			n.child = append(n.child, n1)
		} else {
			// Node already exists
			if h != nil {
				n.handlerFunc = h
				n.paramname = paramname
			}
		}
		return
	}
}
func (r *Router) insert(method, path string, h HandlerFunc, paramname string) {
	n := r.tree[method]
	if n == nil {
		n = &node{}
	}
	n.Add(method, path, h, paramname)
	r.tree[method] = n
}

// func (n *node) addChild(c *node) {
// 	n.children = append(n.children, c)
// }

// func (n *node) findChild(l byte, t kind) *node {
// 	for _, c := range n.children {
// 		if c.label == l && c.kind == t {
// 			return c
// 		}
// 	}
// 	return nil
// }

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.child {
		if c.label == l {
			return c
		}
	}
	return nil
}

// func (n *node) findChildByKind(t kind) *node {
// 	for _, c := range n.children {
// 		if c.kind == t {
// 			return c
// 		}
// 	}
// 	return nil
// }

// func (n *node) checkMethodNotAllowed() HandlerFunc {
// 	for _, m := range methods {
// 		if h := n.findHandler(m); h != nil {
// 			return MethodNotAllowedHandler
// 		}
// 	}
// 	return NotFoundHandler
// }

func (n *node) Find(path string) (h HandlerFunc) {

	var search = path
	// Search order static > param > any
	for {
		l := 0
		sl := len(search)
		pl := len(n.prefix)

		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == n.prefix[l]; l++ {
		}

		if l == pl {
			// Continue search
			search = search[l:]
		} else {
			// cn = nn
			// search = ns
			// if nk == pkind {
			// 	goto Param
			// } else if nk == akind {
			// 	goto Any
			// }
			// Not found
			return nil
		}

		if search == "" {
			return n.handlerFunc
		}

		// Static node
		if n1 := n.findChildWithLabel(search[0]); n1 != nil {
			n = n1
			continue
		}
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, re *http.Request) {
	method := re.Method
	uri := re.URL.Path
	log.Printf("%s %s", method, uri)
	tree, ok := r.tree[method]
	if !ok {
		log.Println("method not allow")
	}
	h := tree.Find(uri)
	if h == nil {
		http.NotFound(rw, re)
		return
	}
	ctx := &Context{
		re,
		rw,
	}
	h(ctx)
}
