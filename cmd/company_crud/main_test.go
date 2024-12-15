package main

import (
	"company-crud/internal/app"
	"company-crud/pkg/http_server"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"company-crud/pkg/producer"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestIntegration_main(t *testing.T) {
	suite := New(t, time.Second*5, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.PwN9mqs6JDOROs42oqojiJ0iGEzOtLejuVrDPITuxqw")

	p, err := suite.postgresContainer.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	host, err := suite.postgresContainer.Host(context.Background())
	require.NoError(t, err)

	kc := initKafka(t, context.Background())
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(kc); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	})

	cfg, err := LoadConfig("../../")
	require.NoError(t, err)

	//logger := zap.NewNop()
	logger, err := logger.New("dev")
	require.NoError(t, err)

	postgresCli, err := postres.New(postres.Config{
		DBHost:     host,
		DBPort:     p.Int(),
		DBUsername: cfg.DBUsername,
		DBPassword: cfg.DBPassword,
		DBName:     cfg.DBName,
	})
	require.NoError(t, err)

	httpServer := http_server.New(cfg.Swagger, cfg.Cors)

	kafkaProducer, err := producer.New(producer.Config{
		Server: "localhost:9092",
		Acks:   cfg.KafkaAcks,
		Topic:  cfg.KafkaTopic,
	})
	require.NoError(t, err)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	defer stop()

	companyCrud := app.New(ctx, logger, app.Config{
		TokenSignature: cfg.TokenSignature,
	}, httpServer, postgresCli, kafkaProducer)

	go companyCrud.Run()

	suite.testCreateHttpCases(t, postgresCli, logger)

}
