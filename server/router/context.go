package router

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	. "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"

	"github.com/julienschmidt/httprouter"
)

type Context struct {
	rw       http.ResponseWriter
	req      *http.Request
	p        httprouter.Params
	handlers []HandlerFunc
	i        int8
	User     *models.User
}

func NewContext(rw http.ResponseWriter,
	req *http.Request,
	p httprouter.Params,
	handlers []HandlerFunc) *Context {
	return &Context{rw, req, p, handlers, -1, nil}
}

// Next executes the next handler in the chain
func (ctx *Context) Next() {
	log.Println(ctx.i)
	ctx.i++
	ctx.handlers[ctx.i](ctx)
}

// P returns a path parameter by name
func (ctx *Context) P(key string) string {
	return ctx.p.ByName(key)
}

// Bind binds a JSON request to a Struct
func (ctx *Context) Bind(obj interface{}) error {
	dec := json.NewDecoder(ctx.req.Body)
	if err := dec.Decode(&obj); err != nil {
		return err
	}
	return nil
}

// Render renders a JSON response
func (ctx *Context) Render(code int, obj interface{}) {
	rw := ctx.rw
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	enc := json.NewEncoder(rw)
	if err := enc.Encode(obj); err != nil {
		log.Panic(err)
	}
}

// Render500 renders a 500 error response using Context.Render
func (ctx *Context) Render500() {
	ctx.Render(http.StatusInternalServerError, NewHttpError(http.StatusInternalServerError, M["error500"]))
}

// RenderOk renders a 200 success response using using Context.Render
func (ctx *Context) RenderOk(msg string) {
	ctx.Render(http.StatusOK, map[string]string{
		"message": msg,
	})
}
