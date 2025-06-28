package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samualhalder/go-restapis/internals/config"
	"github.com/samualhalder/go-restapis/internals/database/sqlite"

	"github.com/samualhalder/go-restapis/internals/http/students"
)

func main() {

	cfg := config.MustLoad()

	_, dberr := sqlite.New(cfg)
	if dberr != nil {
		log.Fatal(dberr)
	}
	slog.Info("Database conectead", slog.String("databate type: ", "sqlite"))

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New())

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server is started at port:", slog.String("", cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Server not started", err.Error())
		}
	}()
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Error while stuhing down the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
