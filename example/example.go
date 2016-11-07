package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"virgin"
)

func test(ctx *virgin.Context) {
	fmt.Fprintf(ctx.Response, "Request Method:%s", ctx.Request.Method)
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
	v.AddRoute("GET", "/ad/:a", test)
	// v.AddRoute("GET", "/abc/*", test)
	// v.AddRoute("GET", "/adf", test)
	// v.AddRoute("GET", "/asd", test)
	// v.AddRoute("GET", "/b", test)
	// v.AddRoute("GET", "/", test)
	// v.AddRoute("POST", "/d", test)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log.Fatal(v.Listen(":8080"))
}
