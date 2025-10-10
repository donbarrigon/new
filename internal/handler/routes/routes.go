package routes

import "donbarrigon/new/internal/utils/handler"

func AppRoutes() *handler.Handler {
	h := handler.New()
	h.Get("/", func(c *handler.Context) {
		c.ResponseOk(map[string]string{"status": "ok"})
	})
	return h
}
