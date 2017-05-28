package virgin

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

var Recovey = func(context *Context) {
	if err := recover(); err != nil {
		if context.Response.Header().Get("Content-Type") == "" {
			context.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}

		context.Response.WriteHeader(http.StatusInternalServerError)
		stack := debug.Stack()
		fmt.Printf("recover in %s,%s", err, string(stack))
	}
}
