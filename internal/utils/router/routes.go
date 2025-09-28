package router

import (
	"donbarrigon/new/internal/utils/handler"
	"net/http"
	"strings"
)

type ControllerFun func(ctx *handler.HttpContext)
type MiddlewareFun func(func(ctx *handler.HttpContext)) func(ctx *handler.HttpContext)

type Route struct {
	Path       []string
	IsVar      []bool
	Controller ControllerFun
	Middleware []MiddlewareFun
	Name       string
}

var (
	routes      = []Route{}
	prefixes    = []string{}
	middlewares = []MiddlewareFun{}
)

func Get(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodGet, path, ctrl, middlewares...)
}

func Post(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodPost, path, ctrl, middlewares...)
}

func Patch(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodPatch, path, ctrl, middlewares...)
}

func Put(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodPut, path, ctrl, middlewares...)
}

func Delete(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodDelete, path, ctrl, middlewares...)
}

func Options(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodOptions, path, ctrl, middlewares...)
}

func Head(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	SetRoute(http.MethodHead, path, ctrl, middlewares...)
}

func Prefix(prefix string, callback func(), mws ...MiddlewareFun) {
	pfs := strings.Split(strings.Trim(prefix, "/"), "/")

	prefixes = append(prefixes, pfs...)
	middlewares = append(middlewares, mws...)
	callback()
	prefixes = prefixes[:len(prefixes)-len(pfs)]
	middlewares = middlewares[:len(middlewares)-len(mws)]
}

func Middleware(callback func(), mws ...MiddlewareFun) {
	middlewares = append(middlewares, mws...)
	callback()
	middlewares = middlewares[:len(middlewares)-len(mws)]
}

func Name(name string) {
	prefix := strings.Join(prefixes, ".")
	routes[len(routes)-1].Name = prefix + "." + name
}

func SetRoute(method string, path string, ctrl ControllerFun, mws ...MiddlewareFun) {
	segments := prefixes
	segments = append(segments, strings.Split(strings.Trim(path, "/"), "/")...)

	var pathParts []string
	var isVars []bool

	pathParts = append(pathParts, method)
	isVars = append(isVars, false)

	for _, part := range segments {
		if strings.HasPrefix(part, ":") || (strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")) {
			part := strings.Trim(part, ":{}")
			pathParts = append(pathParts, part)
			isVars = append(isVars, true)
		} else {
			pathParts = append(pathParts, part)
			isVars = append(isVars, false)
		}
	}

	// copio los middlewares por que sino pasa la referencia que se modifica en otros lados y se rompe
	mwCopy := make([]MiddlewareFun, len(middlewares))
	copy(mwCopy, middlewares)

	routes = append(routes, Route{
		Path:       pathParts,
		IsVar:      isVars,
		Controller: ctrl,
		Middleware: append(mwCopy, mws...),
	})
}
