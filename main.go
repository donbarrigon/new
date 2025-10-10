package main

import (
	"donbarrigon/new/internal/handler/routes"
	"donbarrigon/new/internal/utils/server"
)

//go:generate bash -c "mkdir -p internal/ui/view && qtc -dir=internal/ui/pages && mv internal/ui/pages/*.qtpl.go internal/ui/view/"
func main() {

	server.Start(routes.AppRoutes())
}
