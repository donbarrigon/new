package routes

import (
	"donbarrigon/new/internal/database/handler/controller"
	"donbarrigon/new/internal/database/handler/middleware"
	"donbarrigon/new/internal/utils/handler"
)

func DatabaseRoutes(h *handler.Handler) {
	h.Prefix("db", func() {
		h.Get("migrate", controller.Migrate)
		h.Get("migrate/rollback", controller.Rollback)
		h.Get("migrate/reset", controller.Reset)
		h.Get("migrate/refresh", controller.Refresh)
		h.Get("migrate/fresh", controller.Fresh)
		h.Get("seed", controller.Seed)
		h.Get("seed/all", controller.SeedAll)
		h.Get("seed/run", controller.SeedRun)
	}, middleware.Migration)
}
