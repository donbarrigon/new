package app

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func ServerStart(routes *Routes) {
	router := &Router{}
	router.Make(routes)

	var tlsConfig *tls.Config

	if Env.SERVER_HTTPS_ENABLED {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache(Env.SERVER_HTTPS_CERT_PATH), // Carpeta donde guarda los certs
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(strings.Split(Env.SERVER_HOST_WHITE_LIST, ",")...),
		}
		tlsConfig = manager.TLSConfig()
		go func() {
			// Redirigir todo lo que llegue por HTTP a HTTPS
			if er := http.ListenAndServe(":80", manager.HTTPHandler(nil)); er != nil {
				PrintError("🔴💥 Could not start server http port 80: :error", Entry{"error", er.Error()})
			}
		}()
	}

	server := &http.Server{
		Addr:         ":" + Env.SERVER_PORT,
		Handler:      router.HandlerFunction(),
		TLSConfig:    tlsConfig,
		ReadTimeout:  time.Duration(Env.SERVER_READ_TIMEOUT) * time.Second,
		WriteTimeout: time.Duration(Env.SERVER_WRITE_TIMEOUT) * time.Second,
		IdleTimeout:  time.Duration(Env.SERVER_READ_TIMEOUT+Env.SERVER_WRITE_TIMEOUT) * time.Second,
	}

	PrintInfo(`🚀 Server running on :app_url 
  ____   ___  ____  ____  ___  ___  _   _ ____   ___
 / ___| / _ \|  _ \|  _ \|_ _|| __|| \ | |  _ \ / _ \
| |    | | | | |_) | |_) || ||||__ |  \| | | | | | | |
| |___ | |_| |  _ <|  _ < | ||||__ | |\  | |_| | |_| |
 \____(_)___/|_| \_\_| \_\___||___||_| \_|____/ \___/
`, Entry{"app_url", Env.APP_URL})

	// funciona en dev pero en produccion es feo.
	// espera la señal en segundo plano, el bun run dev lo reinicia pero el main se termina y no salen los mensajes
	go HttpServerGracefulShutdown(server)
	time.Sleep(100 * time.Millisecond) // para que salga el mensaje de corriendo.

	if Env.SERVER_HTTPS_ENABLED {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			PrintError("🔴💥 Could not start server tls: :error", Entry{"error", err.Error()})
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			PrintError("🔴💥 Could not start server: :error", Entry{"error", err.Error()})
		}
	}

	// funciona mejor en produccion en dev no.
	// server en segundo plano espera la señal de cierre, salen los mensajes pero el bun run dev no lo reinicia
	// go func() {
	// 	if Env.SERVER_HTTPS_ENABLED {
	// 		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
	// 			PrintError("🔴💥 Could not start server: :error", Entry{"error", err.Error()})
	// 		}
	// 	} else {
	// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 			PrintError("🔴💥 Could not start server: :error", Entry{"error", err.Error()})
	// 		}
	// 	}
	// }()
	// HttpServerGracefulShutdown(server)
}

// maneja el apagado graceful del servidor
func HttpServerGracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Espera por la señal de terminación
	<-sigChan
	PrintInfo("⏻ Initiating controlled server shutdown...")

	// se cierra la conexion con mono db
	if err := CloseMongoDB(); err != nil {
		PrintWarning("🔴💥 Error closing connection to MongoDB :err", Entry{"err", err.Error()})
	} else {
		PrintInfo("🔌 Connection to MongoDB successfully closed")
	}

	// Crea un contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//se cierra el servidor HTTP para que no acepte nuevas conexiones
	if err := server.Shutdown(ctx); err != nil {
		PrintWarning("⏻ Server forced to close: :err", Entry{"err", err.Error()})
	} else {
		PrintInfo("⏻ HTTP server stopped successfully")
	}

	PrintInfo("💀 Apagado controlado completado")
}
