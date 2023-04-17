package gee

import (
	"log"
	"time"
)


// Logger 日志记录时间中间件, 可以记录每次请求
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
