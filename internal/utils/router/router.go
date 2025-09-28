package router

import (
	"donbarrigon/new/internal/utils/handler"
	"net/http"
	"strings"
)

type Router struct {
	path        string
	isVar       bool
	index       int
	controller  ControllerFun
	middlewares []MiddlewareFun
	routers     []*Router
}
type DinamicRoute struct {
	isVar  bool
	index  int
	routes map[string]DinamicRoute
}

type RouterData struct {
	Params      map[string]string
	Controller  ControllerFun
	Middlewares []MiddlewareFun
}

var (
	StaticRoutes  map[string]ControllerFun
	DinamicRoutes map[string]DinamicRoute
	NameRoutes    map[string]string
)

func Make() {
	StaticRoutes = map[string]ControllerFun{}
	DinamicRoutes = map[string]DinamicRoute{}
	NameRoutes = map[string]string{}
	for _, route := range routes {
		isStatic := true
		for _, isVar := range route.IsVar {
			if isVar {
				isStatic = false
				break
			}
		}
		if isStatic {
			StaticRoutes[strings.Join(route.Path, "/")] = Use(route.Controller, route.Middleware...)
		} else {
			// continuar aca
		}
	}
}

// construlle las rutas optimizadas para luego buscar
// toma el array de rutas y las convierte en ramas
// func (r *Router) Make(routes *Routes) {
// 	r.routers = []*Router{}
// 	// recorro las rutas
// 	for _, route := range routes.routes {
// 		// paso el router padre y la ruta, el indice es cero por que es la raiz
// 		r.add(0, route)
// 	}
// }

// avansa recursivamente por el array del pat de la ruta y va creando ramas
func (r *Router) add(index int, route *Route) {
	// veo a ver si el path ya existe o es una variable existente
	for i, router := range r.routers {
		// if router.isVar == route.IsVar[index] || router.path == route.Path[index] {
		if router.isVar || router.path == route.Path[index] {
			// si existe, tiro palante con el router que ya existente
			r.routers[i].add(index+1, route)
			return
		}
	}

	// si no existe, creo un nuevo router
	newRouter := &Router{
		path:    route.Path[index],
		isVar:   route.IsVar[index],
		index:   index,
		routers: []*Router{}, // creo el array de rutas vacio para evitar errores en la recursividad
	}

	//agrego el router a la lista del router actual(el padre)
	r.routers = append(r.routers, newRouter)

	// si la ruta aun no termina, sigo adelante con el nuevo router
	nextIndex := index + 1
	if nextIndex < len(route.Path) {
		newRouter.add(nextIndex, route)
		return
	}

	// si la ruta termina, agrego el controlador
	newRouter.controller = route.Controller
	newRouter.middlewares = route.Middleware
}

// func (r *Router) Find(path string, rd *RouterData) {
// 	p := strings.Split(strings.Trim(path, "/"), "/")
// 	r.find(p, rd)
// }

// busca recursivamente por las ramas del router y si encuentra el controlador lo asigna en rd
// si el controlador de rd es nil, significa que no encontro la ruta y se debe manejar como 404
func (r *Router) Find(path []string, rd *RouterData) {
	for _, router := range r.routers {
		// Log.Print(router.path)
		if router.isVar || path[router.index] == router.path {
			// si es variable se guarda
			if router.isVar {
				rd.Params[router.path] = path[router.index]
			}

			// Log.Print(">>" + router.path + "<<")

			// si el path aun no termina, sigue adelante y si el siguiente no tiene rutas, se para
			if len(path) > router.index+1 {
				router.Find(path, rd)
				return
			}

			// si es la ultima rama devuelve controlador
			if router.controller != nil {
				rd.Controller = router.controller
				rd.Middlewares = router.middlewares
				return
			}
		}
	}
}

func (router *Router) HandlerFunction() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := handler.NewHttpContext(w, r)

		rd := &RouterData{
			Params: map[string]string{},
		}
		pathSegments := []string{r.Method}
		pathSegments = append(pathSegments, strings.Split(strings.Trim(r.URL.Path, "/"), "/")...)
		router.Find(pathSegments, rd)

		if rd.Controller != nil {
			ctx.Params = rd.Params
			router.Use(rd.Controller, rd.Middlewares...)(ctx)
			return
		}
		ctx.ResponseNotFound()
	}
}

func Use(function ControllerFun, middlewares ...MiddlewareFun) ControllerFun {
	for i := len(middlewares) - 1; i >= 0; i-- {
		function = middlewares[i](function)
	}
	return function
}
