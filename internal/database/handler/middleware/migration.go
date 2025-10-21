package middleware

import (
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/logs"
	"net"
)

func Migration(next handler.ControllerFun) handler.ControllerFun {
	return func(c *handler.Context) {
		if !config.DbMigrationEnable {
			logs.Warning("Migraciones deshabilitadas, Habilitelas en el archivo .env DB_MIGRATION_ENABLE=true")
			c.ResponseNoContent()
			return
		}

		host, _, er := net.SplitHostPort(c.Request.RemoteAddr)
		if er != nil {
			logs.Warning("Error al obtener la direccion remota: " + er.Error())
			c.ResponseNoContent()
			return
		}
		if host == "127.0.0.1" || host == "::1" {
			next(c)
			return
		}
		logs.Warning("La direccion remota no es localhost")
		c.ResponseNoContent()
	}
}
