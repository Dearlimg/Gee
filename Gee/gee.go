package gee

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// type HandlerFunc func(http.ResponseWriter, *http.Request)
// 函数签名
type HandlerFunc func(c *Context)

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func New() *Engine {
	//return &Engine{router: newRouter()}
	//-----------------------------------
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//	engine.router.addRoute(method, pattern, handler)
//}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//	engine.addRoute("GET", pattern, handler)
//}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//	engine.addRoute("POST", pattern, handler)
//}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//key := req.Method + "-" + req.URL.Path
	//if handler, ok := engine.router[key]; ok {
	//	handler(w, req)
	//} else {
	//	fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	//}
	//------------------------------
	//c := newContext(w, req)
	//engine.router.handle(c)
	//-------------------------------
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
