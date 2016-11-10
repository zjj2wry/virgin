package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"virgin"
)

func test(ctx *virgin.Context) {
	userid:=ctx.Param("a")
	fmt.Fprintf(ctx.Response, "Request Method:%s,paramvalue:%s", ctx.Request.Method,userid)
}
func test2(ctx *virgin.Context) {
	fmt.Fprintf(ctx.Response, "Request:%s", ctx.Request.Method)
}

// func test3(ctx *virgin.Context) {
// 	fmt.Fprintf(ctx.Response, "Request")
// }

func main() {
	v := virgin.NewVirgin()
	v.AddRoute("GET", "/a", test)
	v.AddRoute("GET", "/ab", test2)
	v.AddRoute("GET", "add", test)
	// v.AddRoute("GET", "/b/*", test)
	v.AddRoute("GET", "/c/:a", test)
	// v.AddRoute("GET", "/adf", test)
	// v.AddRoute("GET", "/asd", test)
	// v.AddRoute("GET", "/b", test)
	// v.AddRoute("GET", "/", test)
	// v.AddRoute("POST", "/d", test)
	// fmt.Println(v.Router)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Fatal(v.Listen(":8080"))
}
