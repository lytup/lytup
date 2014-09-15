package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type (
	HandlerFunc func(*Context)

	Group struct {
		mux      *Mux
		prefix   string
		handlers []HandlerFunc
	}

	Mux struct {
		router *httprouter.Router
		*Group
	}
)

// New creates a new router.
func New() *Mux {
	r := &Mux{}
	r.router = httprouter.New()
	r.Group = &Group{r, "", nil}
	return r
}

func (m *Mux) Run() {
	http.ListenAndServe("localhost:3000", m.router)
}

// G creates a new group with prefix and handlers
// Root handlers are prepended to the chain
func (g *Group) G(prefix string, handlers ...HandlerFunc) *Group {
	return &Group{
		g.mux,
		prefix,
		append(g.mux.handlers, handlers...), // Prepend root handlers
	}
}

func (g *Group) Handle(method, path string, handlers []HandlerFunc) {
	handlers = append(g.handlers, handlers...) // Prepend group handlers
	path = g.prefix + path
	g.mux.router.Handle(method, path, func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx := NewContext(rw, req, p, handlers)
		ctx.Next()
	})
}

// Use adds new middleware(s)
func (g *Group) Use(middlewares ...HandlerFunc) {
	g.handlers = append(g.handlers, middlewares...)
}

func (g *Group) Post(path string, handlers ...HandlerFunc) {
	g.Handle("POST", path, handlers)
}

func (g *Group) Get(path string, handlers ...HandlerFunc) {
	g.Handle("GET", path, handlers)
}

func (g *Group) Delete(path string, handlers ...HandlerFunc) {
	g.Handle("DELETE", path, handlers)
}

func (g *Group) Patch(path string, handlers ...HandlerFunc) {
	g.Handle("PATCH", path, handlers)
}

func (g *Group) Put(path string, handlers ...HandlerFunc) {
	g.Handle("PUT", path, handlers)
}
