package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahulkumarpahwa/module-manager/config"
	"github.com/rahulkumarpahwa/module-manager/internal/handlers"
	"github.com/rahulkumarpahwa/module-manager/internal/middleware"
	"github.com/rahulkumarpahwa/module-manager/services"
	"github.com/rahulkumarpahwa/module-manager/services/types"
)

func main() {
	log.SetFlags(log.Ltime | log.Ldate)

	cfg, err := config.Open()
	if err != nil {
		log.Fatalf("main.config.Open: %v", err.Error())
	}

	mainFolderInfo := types.NewMainFolder()
	services := services.NewServices(mainFolderInfo)
	handlers := handlers.NewHandlers(mainFolderInfo, services)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handlers.Health)
	mux.HandleFunc("POST /api/set-path", handlers.SetPath)
	mux.HandleFunc("POST /api/main/scan", handlers.ScanMainFolder)
	mux.HandleFunc("POST /api/project/open-vscode", handlers.OpenVSCode)
	mux.HandleFunc("POST /api/project/open-winexplorer", handlers.OpenWinExplorer)
	mux.HandleFunc("GET /api/project/{id}", handlers.GetProjectInfoById)

	loggingMux := middleware.LoggingMiddleware(mux)

	corsMux := middleware.CorsMiddleware(loggingMux)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: corsMux,
		//todo: task set these values in the config as well below.
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server is Listening at http://localhost:%s", cfg.Port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				// skip
			} else {
				log.Fatalf("main.server.ListenAndServer: %v", err)
			}
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Print("Server Shutdown in...")
	for i := 5; i > 0; i-- {
		log.Printf("%d", i)
		time.Sleep(time.Second * 1)
	}

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server.shutdown: %v", err)
	}

	log.Print("server shutdown successfully.")
}
