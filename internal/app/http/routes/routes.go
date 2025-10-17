package routes

import (
	"donbarrigon/new/internal/app/http/controller"
	"donbarrigon/new/internal/utils/handler"
)

func AppRoutes() *handler.Handler {
	h := handler.New()
	h.Get("/", controller.Home)
	return h
}
