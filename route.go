package virgin

import (
	"log"
	"net/http"
	"fmt"
	"time"
)

const (
	NOTFOUND string ="NOT FOUND"
	FOUND	 string = "FOUND"
)

type (
	HandlerFunc func(*Context)

	Router struct {
		tree map[string]*node
	}
	node struct {
		label       byte
		prefix      string
		child       []*node
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
				r.insert(method, path[:i], nil)
				r.insert(method, path, h)
				fmt.Println(path[:i],path)
				return
			}
		} else if path[i] == '*' {
			if path[i-1] == '/' {
				r.insert(method, path[:i], nil)
				r.insert(method, path, h)
				return
			}
		}
	}

	r.insert(method, path, h)
}

func (n *node) Add(path string, h HandlerFunc) {
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
			}
		} else if l < pl {
			n1 := &node{
				n.label,
				n.prefix,
				n.child,
				n.handlerFunc,
			}

			n.label = n.prefix[0]
			n.prefix = n.prefix[:l]

			n.child = append(n.child, n1)

			if l == sl {
				n.handlerFunc = h

			} else {
				prefix := search[l:]
				n2 := &node{
					prefix[0],
					prefix,
					nil,
					h,
				}
				n.child = append(n.child, n2)
			}
		} else if l < sl {
			search = search[l:]
			c := n.findChildWithLabel(search[0])
			if c != nil {
				// 继续检索
				n = c
				continue
			}
			n1 := &node{
				search[0],
				search,
				nil,

				h,
			}
			n.child = append(n.child, n1)
		} else {
			// 遇到参数路由添加的nil handle
			if h != nil {
				n.handlerFunc = h
			}
		}
		return
	}
}

func (r *Router) insert(method, path string, h HandlerFunc) {
	n := r.tree[method]
	if n == nil {
		n = &node{}
	}
	n.Add(path, h)
	r.tree[method] = n
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.child {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) Find(path string) (*node, string) {
	var search = path

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
			search = search[l:]
		} else {
			// 参数路由
			for _, v := range n.child {
				if v.label == ':' {
					return v, search
				}
			}
			//全部匹配
			for _, v := range n.child {
				if v.label == '*' {
					return v, ""
				}
			}
			return nil, ""
		}
		//绝对路由
		if search == "" {
			return n, ""
		}

		if n1 := n.findChildWithLabel(search[0]); n1 != nil {
			n = n1
			continue
		}
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, re *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	t:=time.Now()
	method := re.Method
	uri := re.URL.Path

	tree, ok := r.tree[method]
	
	if !ok {
		http.NotFound(rw,re)
		return
	}
	n, paramname := tree.Find(uri)
	if n == nil|| n.handlerFunc == nil {
		http.NotFound(rw, re)
		dur:=time.Since(t)
		log.Printf("%s %10s %10s %10s", method, uri,dur.String(),NOTFOUND)
		return
	}
	ctx := &Context{
		Request:  re,
		Response: rw,
	}
	// 设置参数获取
	if paramname != "" {
		ctx.setParamname(n.prefix[1:])
		ctx.setParamvalue(paramname)	
	}
	n.handlerFunc(ctx)
	dur:=time.Since(t)
	log.Printf("\033[32m%s %10s %10s %10s", method, uri,dur.String(),FOUND)
}                         
