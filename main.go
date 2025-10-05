package main

import (
	"donbarrigon/new/internal/routes"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/server"
)

func main() {
	config.LoadEnv()
	db.InitMongoDB()
	server.Start(routes.AppRoutes())
}
