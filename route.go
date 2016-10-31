package virgin

import (
	"fmt"
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

// func (r *Router) Find(method, path string, context Context) {
// 	cn := r.tree // Current node as root

// 	var (
// 		search  = path
// 		c       *node  // Child node
// 		n       int    // Param counter
// 		nk      kind   // Next kind
// 		nn      *node  // Next node
// 		ns      string // Next search
// 		pvalues = context.ParamValues()
// 	)

// 	// Search order static > param > any
// 	for {
// 		if search == "" {
// 			goto End
// 		}

// 		pl := 0 // Prefix length
// 		l := 0  // LCP length

// 		if cn.label != ':' {
// 			sl := len(search)
// 			pl = len(cn.prefix)

// 			// LCP
// 			max := pl
// 			if sl < max {
// 				max = sl
// 			}
// 			for ; l < max && search[l] == cn.prefix[l]; l++ {
// 			}
// 		}

// 		if l == pl {
// 			// Continue search
// 			search = search[l:]
// 		} else {
// 			cn = nn
// 			search = ns
// 			if nk == pkind {
// 				goto Param
// 			} else if nk == akind {
// 				goto Any
// 			}
// 			// Not found
// 			return
// 		}

// 		if search == "" {
// 			goto End
// 		}

// 		// Static node
// 		if c = cn.findChild(search[0], skind); c != nil {
// 			// Save next
// 			if cn.prefix[len(cn.prefix)-1] == '/' { // Issue #623
// 				nk = pkind
// 				nn = cn
// 				ns = search
// 			}
// 			cn = c
// 			continue
// 		}

// 		// Param node
// 	Param:
// 		if c = cn.findChildByKind(pkind); c != nil {
// 			// Issue #378
// 			if len(pvalues) == n {
// 				continue
// 			}

// 			// Save next
// 			if cn.prefix[len(cn.prefix)-1] == '/' { // Issue #623
// 				nk = akind
// 				nn = cn
// 				ns = search
// 			}

// 			cn = c
// 			i, l := 0, len(search)
// 			for ; i < l && search[i] != '/'; i++ {
// 			}
// 			pvalues[n] = search[:i]
// 			n++
// 			search = search[i:]
// 			continue
// 		}

// 		// Any node
// 	Any:
// 		if cn = cn.findChildByKind(akind); cn == nil {
// 			if nn != nil {
// 				cn = nn
// 				nn = nil // Next
// 				search = ns
// 				if nk == pkind {
// 					goto Param
// 				} else if nk == akind {
// 					goto Any
// 				}
// 			}
// 			// Not found
// 			return
// 		}
// 		pvalues[len(cn.paramname)-1] = search
// 		goto End
// 	}

// End:
// 	context.SetHandler(cn.findHandler(method))
// 	context.SetPath(cn.ppath)
// 	context.SetParamNames(cn.paramname...)

// 	// NOTE: Slow zone...
// 	if context.Handler() == nil {
// 		context.SetHandler(cn.checkMethodNotAllowed())

// 		// Dig further for any, might have an empty value for *, e.g.
// 		// serving a directory. Issue #207.
// 		if cn = cn.findChildByKind(akind); cn == nil {
// 			return
// 		}
// 		if h := cn.findHandler(method); h != nil {
// 			context.SetHandler(h)
// 		} else {
// 			context.SetHandler(cn.checkMethodNotAllowed())
// 		}
// 		context.SetPath(cn.ppath)
// 		context.SetParamNames(cn.paramname...)
// 		pvalues[len(cn.paramname)-1] = ""
// 	}

// 	return
// }

func (r *Router) ServeHTTP(rw http.ResponseWriter, re *http.Request) {
	method := re.Method
	uri := re.URL.Path
	log.Printf("%s %s", method, uri)
	handleMap := r.tree["GET"]
	fmt.Println(r.tree)
	fmt.Println(r.tree["GET"])
	handfunc := handleMap.handlerFunc
	ctx := &Context{
		re,
		rw,
	}
	handfunc(ctx)
}
