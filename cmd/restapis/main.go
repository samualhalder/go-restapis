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
)

func main() {

	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to go rest api"))
	})

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server is started ", slog.String("at port", cfg.Addr))

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
