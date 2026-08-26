package main

import (
	"context"
	"log"
	"module-manager/config"
	"module-manager/internal/handlers"
	"module-manager/services"
	"module-manager/types"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Ltime | log.Ldate)

	cfg, err := config.Open()
	if err != nil {
		log.Fatalf("main.config.Open: %v", err.Error())
	}

	mainFolderInfo := types.NewMainFolderInfo()
	services := services.NewServices(mainFolderInfo)
	handlers := handlers.NewHandlers(mainFolderInfo, services)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("POST /api/set-path", handlers.SetPath)
	mux.HandleFunc("POST /api/folders/scan", handlers.ScanFolder)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		// task set these values in the config as well below.
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server is Listening at http://localhost:%s", cfg.Port)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
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
		time.Sleep(time.Millisecond * 10)
	}

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server.shutdown: %v", err)
	}

	log.Print("server shutdown successfully.")
}
