package virgin

import (
	"net/http"
)

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	paramname  string
	paramvalue string
}

func (c Context) setParamname(s string) {
	c.paramname = s
}

func (c Context) setParamvalue(i string) {
	c.paramvalue = i
}

func(c Context) Param(key string)(value string){
	if key==c.paramname{
		return c.paramvalue
	}
	return ""
}