package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"backend-boiler/internal/app"
	"backend-boiler/internal/infra/config"
)

func main() {
	cfg := config.Load()
	container, err := app.NewContainer(cfg)
	if err != nil {
		slog.Error("failed to build container", "error", err)
		os.Exit(1)
	}
	defer container.Close(context.Background())

	worker := app.BootstrapWorker(container)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	container.Logger.Info("worker stopped")
}
