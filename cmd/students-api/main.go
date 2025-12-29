package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OmarAkhter/golan-project/internal/config"
	"github.com/OmarAkhter/golan-project/internal/http/hanlders/students"
)

func main() {
	// loads config from file or environment variables
	cfg := config.MustLoad()

	// route setup
	router := http.NewServeMux()

	router.HandleFunc("POST /", students.New())

	//server setup

	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Starting Server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server: ", err)
		}

	}()

	<-done

	slog.Info("Server Shutting Down")

	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to stop server: ", slog.String("error", err.Error()))
	}

	slog.Info("Server Exited Properly")

}
