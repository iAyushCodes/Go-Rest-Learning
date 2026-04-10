package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iAyushCodes/Go-Rest-Learning/internal/config"
	"github.com/iAyushCodes/Go-Rest-Learning/internal/http/handlers/student"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// db setup

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server started: %s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // these signals will be delivered to 'done' channel

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done // initially we will be blocked here but when 'done' receives the signal, we will be unblocked and the code below this, will start executing

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
