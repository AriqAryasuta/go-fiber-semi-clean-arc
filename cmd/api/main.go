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

	fiberApp := app.BootstrapAPI(container)
	go func() {
		if err := fiberApp.Listen(cfg.ServerAddress()); err != nil {
			container.Logger.Error("api server stopped", "error", err)
		}
	}()

	container.Logger.Info("api server started", "address", cfg.ServerAddress())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
		container.Logger.Error("failed shutting down api server", "error", err)
	}
}
