package main

import (
	"fmt"
	"log"
	"virgin"
)

func test(ctx *virgin.Context) {
	fmt.Fprintf(ctx.Response, "Request Method:%s", ctx.Request.URL.PATH)
}

func main() {
	v := virgin.NewVirgin()
	v.AddRoute("GET", "/a", test)

	log.Fatal(v.Listen(":8080"))
}
