package app

import (
	"company-crud/internal/handlers/http"
	"company-crud/internal/repositories/db"
	"company-crud/internal/services"
	"company-crud/pkg/http_server"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"company-crud/pkg/producer"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	TokenSignature string
}

type CompanyCRUD struct {
	osSignalContext context.Context
	log             *logger.Logger
	server          *http_server.Server
	db              *postres.Postgres
	producer        *producer.KafkaProducer
	cfg             Config
}

func New(ctx context.Context, log *logger.Logger, cfg Config, server *http_server.Server, db *postres.Postgres, producer *producer.KafkaProducer) *CompanyCRUD {
	return &CompanyCRUD{
		osSignalContext: ctx,
		log:             log,
		server:          server,
		db:              db,
		producer:        producer,
		cfg:             cfg,
	}
}

func (cc *CompanyCRUD) Run() {
	cc.log.Info("Company CRUD started...")

	companyDB := db.New(cc.db, cc.log)
	companyService := services.New(cc.log, cc.producer, companyDB)
	companyHttp := http.New(cc.log, companyService, cc.cfg.TokenSignature)

	cc.server.CreateRoutes(companyHttp)
	go func() {
		cc.log.Info(fmt.Sprintf("Listening on: %s", "8000"))
		err := cc.server.Start()
		if err != nil {
			cc.log.Fatal(fmt.Sprintf("error on server %v", err))
		}
	}()

	select {
	case <-cc.osSignalContext.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cc.log.Info("shutting down gracefully...")

		cc.log.Info("shutting down Server...")
		if err := cc.server.Stop(ctx); err != nil {
			cc.log.Fatal("server shutdown failed: %v", zap.Error(err))
		}

		cc.log.Info("shutting down db...")
		if err := cc.db.Stop(); err != nil {
			cc.log.Fatal("db shutdown failed: %v", zap.Error(err))
		}

		cc.log.Info("shutting down kafka producer...")
		cc.producer.Stop()
	}
}
