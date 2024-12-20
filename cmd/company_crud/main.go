package main

import (
	_ "company-crud/api" //swagger
	"company-crud/internal/app"
	"company-crud/pkg/http_server"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"company-crud/pkg/producer"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title CompanyCrud
// @version 0.1
// @description APIs for a company handler - JWT Token without expiration date: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.PwN9mqs6JDOROs42oqojiJ0iGEzOtLejuVrDPITuxqw

// @contact.name b10z

// @license.name None

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Token

// @host localhost:8000

func main() {
	cfg, err := LoadConfig("../")
	if err != nil {
		slog.Error("init failed:", "error", err)
		os.Exit(0)
	}

	log, err := logger.New(cfg.Environment)
	if err != nil {
		slog.Error("logger init failed:", "error", err)
		os.Exit(0)
	}

	postgres, err := postres.New(postres.Config{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBUsername: cfg.DBUsername,
		DBPassword: cfg.DBPassword,
		DBName:     cfg.DBName,
	})
	if err != nil {
		slog.Error("pg init failed:", "error", err)
		os.Exit(0)
	}

	httpServer := http_server.New(cfg.Swagger, cfg.Cors)

	kafkaProducer, err := producer.New(producer.Config{
		Server: cfg.KafkaServer,
		Acks:   cfg.KafkaAcks,
		Topic:  cfg.KafkaTopic,
	})
	if err != nil {
		slog.Error("kafka producer init failed:", "error", err)
		os.Exit(0)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	defer stop()

	companyCrud := app.New(ctx, log, app.Config{
		TokenSignature: cfg.TokenSignature,
	}, httpServer, postgres, kafkaProducer)

	companyCrud.Run()
}
