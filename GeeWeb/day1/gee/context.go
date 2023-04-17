package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
上下文
封装 request response，提供对json，html等返回类型对支持
*/

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response code
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int // 记录当前执行到第几个中间件
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	//v1
	//return &Context{
	//	Writer: w,
	//	Req:    req,
	//	Path:   req.URL.Path,
	//	Method: req.Method,
	//}

	// v2 中间件
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next 遍历执行handler list
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c) // 执行对应对路由方法（中间件也在此处遍历被执行）
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

//post表单
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//get query
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, val string) {
	c.Writer.Header().Set(key, val)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	val, _ := c.Params[key]
	return val
}
