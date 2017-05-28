package virgin

import (
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	paramname  string
	paramvalue string
}

func (c *Context) setParamname(s string) {
	c.paramname = s
}

func (c *Context) setParamvalue(s string) {
	c.paramvalue = s
}

func (c *Context) Param(key string) (value string) {
	if key == c.paramname {
		return c.paramvalue
	}

	return c.paramvalue
}

func (c *Context) Json(i interface{}) {
	b, err := json.Marshal(&i)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
	c.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Response.Write(b)
}
