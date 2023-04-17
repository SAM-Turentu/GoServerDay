package gee

import (
	"log"
	"net/http"
	"strings"
)

/*

自定义一个web框架

*/

type HandlerFunc func(c *Context) // != http.HandlerFunc

// v3
type (
	RouteGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouteGroup
		e           *Engine
	}

	Engine struct {
		*RouteGroup
		router *router
		groups []*RouteGroup
	}
)

// ▲ 升级上方的结构体
//type Engine struct {
//	//v1
//	//router map[string]http.HandlerFunc
//
//	//v2
//	router *router
//}

func New() *Engine {
	//v1
	//return &Engine{router: make(map[string]http.HandlerFunc)}

	//v2
	//return &Engine{router: newRouter()}

	//v3 组路由
	engine := &Engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{e: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

// Group 创建新的RouteGroup
func (group *RouteGroup) Group(prefix string) *RouteGroup {
	engine := group.e
	newGroup := &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		e:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 定义一个增加组的中间件
func (group *RouteGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

//func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
func (group *RouteGroup) addRoute(method, comp string, handler HandlerFunc) {
	//v1
	//key := method + "-" + pattern
	//e.router[key] = handler

	//v2
	//e.router.addRoute(method, pattern, handler)

	//v3
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.e.router.addRoute(method, pattern, handler)
}

//func (e *Engine) GET(pattern string, handler HandlerFunc) {
func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//func (e *Engine) POST(pattern string, handler HandlerFunc) {
func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//v1
	//key := req.Method + "-" + req.URL.Path
	//if handler, ok := e.router[key]; ok {
	//	handler(w, req)
	//} else {
	//	fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	//}

	//v2
	//c := newContext(w, req)
	//e.router.handle(c)

	//v3
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	e.router.handle(c)

}
