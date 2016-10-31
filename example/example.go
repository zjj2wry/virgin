package main

import (
	"fmt"
	"log"
	"virgin"
)

func test(ctx *virgin.Context) {
	fmt.Fprintf(ctx.Response, "Request Method:%s", ctx.Request.Method)
}

func main() {
	v := virgin.NewVirgin()
	v.AddRoute("GET", "/a", test)
	v.AddRoute("GET", "/ab", test)
	v.AddRoute("GET", "/ad/:a", test)
	v.AddRoute("GET", "/av/", test)
	v.AddRoute("GET", "/adf", test)
	v.AddRoute("GET", "/asd", test)
	v.AddRoute("GET", "/b", test)
	v.AddRoute("GET", "/", test)
	v.AddRoute("POST", "/dsf", test)

	log.Fatal(v.Listen(":8080"))
}
