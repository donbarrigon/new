package routes

import (
	"donbarrigon/new/internal/handler/controller"
	"donbarrigon/new/internal/utils/handler"
)

func AppRoutes() *handler.Handler {
	h := handler.New()
	h.Get("/", controller.Home)
	return h
}
