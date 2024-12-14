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

type CompanyCRUD struct {
	osSignalContext context.Context
	log             *logger.Logger
	server          *http_server.Server
	db              *postres.Postgres
	producer        *producer.KafkaProducer
}

func New(ctx context.Context, log *logger.Logger, server *http_server.Server, db *postres.Postgres, producer *producer.KafkaProducer) *CompanyCRUD {
	return &CompanyCRUD{
		osSignalContext: ctx,
		log:             log,
		server:          server,
		db:              db,
		producer:        producer,
	}
}

func (cc *CompanyCRUD) Run() {
	companyDB := db.New(cc.db, cc.log)
	companyService := services.New(cc.log, cc.producer, companyDB)
	companyHttp := http.New(cc.log, companyService)

	cc.server.CreateRoutes(companyHttp)
	cc.log.Info(fmt.Sprintf("Listening on: %s", "8000"))
	go func() {
		cc.log.Info(fmt.Sprintf("Listening on: %s", "8000"))
		err := cc.server.Start()
		if err != nil {
			cc.log.Fatal(fmt.Sprintf("error on server %v", err))
		}
	}()

	select { //nolint:all
	case <-cc.osSignalContext.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cc.log.Info("shutting down gracefully...")
		cc.log.Info("shutting down Server...")
		if err := cc.server.Stop(ctx); err != nil {
			cc.log.Fatal("server shutdown failed: %v", zap.Error(err))
		}
	}
}
