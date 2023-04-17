package gee

import (
	"net/http"
	"strings"
)

/*

路由

*/

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	//v1
	//log.Printf("Route %4s - %s", method, pattern) // 打印所有路由
	//key := method + "-" + pattern
	//r.handlers[key] = handler

	//v2
	parts := parsePattern(pattern) // 分割pattern
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler  // 将路由和对应对handler方法一一对应，等待在 handle 中执行
}

//func (r *router) handle(c *Context) {
//	//v1
//	//key := c.Method + "-" + c.Path
//	//if handler, ok := r.handlers[key]; ok {
//	//	handler(c)
//	//} else {
//	//	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
//	//}
//
//	//v2
//	n, params := r.getRoute(c.Method, c.Path)
//	if n != nil {
//		c.Params = params
//		key := c.Method + "-" + n.pattern
//		r.handlers[key](c)  // 执行对应对路由方法
//	}else {
//		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
//	}
//}

// handle 带有中间件
func (r *router) handle(c *Context)  {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])  // 将对应的路由方法添加到 handlers 列表中
	}else{
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()  // 调用遍历方法执行handler
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
