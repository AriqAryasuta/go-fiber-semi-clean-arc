package app

import (
	"context"

	"backend-boiler/internal/http/middleware"
	"backend-boiler/internal/http/routes"
	"backend-boiler/internal/platform/broker"
	"github.com/gofiber/fiber/v3"
	"github.com/rabbitmq/amqp091-go"
)

func BootstrapAPI(container *Container) *fiber.App {
	app := fiber.New()

	app.Use(middleware.Recovery())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logging(container.Logger))

	routes.Register(app)
	return app
}

type Worker struct {
	logger   Logger
	consumer *broker.Consumer
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

func BootstrapWorker(container *Container) *Worker {
	consumer, err := broker.NewConsumer(container.RabbitMQ)
	if err != nil {
		container.Logger.Error("failed to init consumer", "error", err)
		return &Worker{logger: container.Logger}
	}
	return &Worker{
		logger:   container.Logger,
		consumer: consumer,
	}
}

func (w *Worker) Start(ctx context.Context) {
	if w.consumer == nil {
		w.logger.Info("worker started without consumer")
		<-ctx.Done()
		return
	}

	w.logger.Info("worker started", "queue", "jobs.default")
	err := w.consumer.Start(ctx, "jobs.default", func(ctx context.Context, msg amqp091.Delivery) error {
		w.logger.Info("message received", "message_id", msg.MessageId)
		_ = ctx
		return nil
	})
	if err != nil {
		w.logger.Error("worker consume stopped", "error", err)
	}
	_ = w.consumer.Close()
}
