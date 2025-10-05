package handler

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/logs"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"
)

type ControllerFun func(ctx *Context)
type MiddlewareFun func(ControllerFun) ControllerFun

type Handler struct {
	// http.Handler
	// prefix/path:method -> controller
	Routes map[string]ControllerFun
	// default method.prefix.path -> route
	RoutesNames map[string]string
	middlewares []MiddlewareFun
	prefixes    []string
	lastRoute   string
	lastName    string
}

func New() *Handler {
	return &Handler{
		Routes:      map[string]ControllerFun{},
		RoutesNames: map[string]string{},
		middlewares: []MiddlewareFun{},
		prefixes:    []string{},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r, h)
	defer func() {
		if er := recover(); er != nil {
			e := err.Panic(er, string(debug.Stack()))
			msg, er2 := json.Marshal(e.Err)
			if er2 != nil {
				msg = []byte{}
			}
			go logs.Critical("Panic: %s", string(msg))
			c.ResponseError(e)
			return
		}
	}()
	route := strings.ToLower(strings.Trim(r.URL.Path, "/")) + ":" + r.Method
	if ctrl, ok := h.Routes[route]; ok {
		ctrl(c)
		return
	}
	c.ResponseNotFound()
}

func (h *Handler) Get(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodGet, path, ctrl, middlewares...)
}

func (h *Handler) Post(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodPost, path, ctrl, middlewares...)
}

func (h *Handler) Put(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodPut, path, ctrl, middlewares...)
}

func (h *Handler) Patch(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodPatch, path, ctrl, middlewares...)
}

func (h *Handler) Delete(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodDelete, path, ctrl, middlewares...)
}

func (h *Handler) Head(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodHead, path, ctrl, middlewares...)
}

func (h *Handler) Options(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodOptions, path, ctrl, middlewares...)
}

func (h *Handler) Connect(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodConnect, path, ctrl, middlewares...)
}

func (h *Handler) Trace(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	h.add(http.MethodTrace, path, ctrl, middlewares...)
}

func (h *Handler) Prefix(prefix string, callback func(), mws ...MiddlewareFun) {
	pfs := strings.Split(strings.Trim(prefix, "/"), "/")

	h.prefixes = append(h.prefixes, pfs...)
	h.middlewares = append(h.middlewares, mws...)
	callback()
	h.prefixes = h.prefixes[:len(h.prefixes)-len(pfs)]
	h.middlewares = h.middlewares[:len(h.middlewares)-len(mws)]
}

func (h *Handler) Middleware(callback func(), mws ...MiddlewareFun) {
	h.middlewares = append(h.middlewares, mws...)
	callback()
	h.middlewares = h.middlewares[:len(h.middlewares)-len(mws)]
}

func (h *Handler) Name(name string) {
	delete(h.RoutesNames, h.lastName)
	name = strings.Join(h.prefixes, ".") + "." + name
	h.RoutesNames[name] = h.lastRoute
}
func (h *Handler) Namef(name string) {
	delete(h.RoutesNames, h.lastName)
	h.RoutesNames[name] = h.lastRoute
}

func (h *Handler) add(method string, path string, ctrl ControllerFun, mws ...MiddlewareFun) {
	prefix := strings.Join(h.prefixes, "/")
	route := strings.ToLower(prefix+"/"+strings.Trim(path, "/")) + ":" + method
	// if route == "" {
	// 	route = "/"
	// }
	mws = append(h.middlewares, mws...)
	mwCopy := make([]MiddlewareFun, len(mws))
	copy(mwCopy, mws)
	for i := len(mwCopy) - 1; i >= 0; i-- {
		ctrl = mwCopy[i](ctrl)
	}
	name := method + "." + strings.Join(h.prefixes, ".") + "." + strings.Replace(strings.Trim(path, "/"), "/", ".", -1)
	name = strings.ToLower(name)
	route = strings.Trim(route, "/")

	h.Routes[route] = ctrl
	h.RoutesNames[name] = route

	h.lastRoute = route
	h.lastName = name
}
