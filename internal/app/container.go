package app

import (
	"context"
	"log/slog"

	"backend-boiler/internal/infra/broker"
	"backend-boiler/internal/infra/cache"
	"backend-boiler/internal/infra/config"
	"backend-boiler/internal/infra/database"
	"backend-boiler/internal/infra/logger"
	sharedValidator "backend-boiler/internal/shared/validator"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	Config    config.Config
	Logger    *slog.Logger
	DB        *gorm.DB
	Redis     *redis.Client
	RabbitMQ  *amqp091.Connection
	Validator *sharedValidator.Validator
}

func NewContainer(cfg config.Config) (*Container, error) {
	log := logger.New(cfg.AppEnv)

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}

	redisClient, err := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return nil, err
	}

	rabbitConn, err := broker.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("failed to connect to rabbitmq", "error", err)
		return nil, err
	}

	slog.Info("connected to database", "url", cfg.DatabaseURL)
	slog.Info("connected to redis", "url", cfg.RedisAddr)
	slog.Info("connected to rabbitmq", "url", cfg.RabbitMQURL)

	return &Container{
		Config:    cfg,
		Logger:    log,
		DB:        db,
		Redis:     redisClient,
		RabbitMQ:  rabbitConn,
		Validator: sharedValidator.New(),
	}, nil
}

func (c *Container) Close(ctx context.Context) {
	if c.Redis != nil {
		_ = c.Redis.Close()
	}
	if c.RabbitMQ != nil {
		_ = c.RabbitMQ.Close()
	}

	sqlDB, err := c.DB.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	_ = ctx
}
