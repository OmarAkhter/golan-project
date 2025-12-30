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
	"github.com/OmarAkhter/golan-project/internal/storage/sqlite"
)

func main() {
	// loads config from file or environment variables
	cfg := config.MustLoad()

	//database setup

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	slog.Info("Storage Initialized")

	// route setup
	router := http.NewServeMux()

	router.HandleFunc("POST /students", students.New(storage))
	router.HandleFunc("GET /students/{id}", students.GetById(storage))
	router.HandleFunc("GET /students", students.GetList(storage))
	router.HandleFunc("DELETE /students/{id}", students.DeleteById(storage))

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

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server Exited Properly")

}
