# virgin

golang web framework

# already implement feature

radix tree route;
any param route;
config read;
recover handle

#### Install

    go get -v github.com/zjj2wry/virgin

#### Usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"virgin"
)

type User struct {
	User string
	Age  int64
}

func test(ctx *virgin.Context) {
	userid := ctx.Param("a")
	panic("baizhi")
	fmt.Fprintf(ctx.Response, "Request Method:%s,paramvalue:%s", ctx.Request.Method, userid)
}
func test3(ctx *virgin.Context) {
	userid := ctx.Param("a")
	fmt.Fprintf(ctx.Response, "Request Method:%s,paramvalue:%s", ctx.Request.Method, userid)
}
func test2(ctx *virgin.Context) {
	ctx.Json(User{
		"baizhi",
		20,
	})
}

func main() {
	v := virgin.NewVirgin()
	v.AddRoute("GET", "/a", test)
	v.AddRoute("POST", "/ab", test2)
	v.AddRoute("GET", "add", test)
	v.AddRoute("GET", "add/*", test3)
	// v.AddRoute("GET", "/b/*", test)
	v.AddRoute("GET", "/c/:a", test)
	// v.AddRoute("GET", "/adf", test)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Fatal(v.Listen(":8080"))
}
```