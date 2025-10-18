package routes

import (
	"donbarrigon/new/internal/database/handler/controller"
	"donbarrigon/new/internal/utils/handler"
)

func DatabaseRoutes(h *handler.Handler) {
	h.Prefix("db", func() {
		h.Get("migrate", controller.Migrate)
	})

}
