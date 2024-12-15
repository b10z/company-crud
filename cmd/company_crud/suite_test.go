package main

import (
	"bytes"
	"company-crud/pkg/logger"
	pkgPg "company-crud/pkg/postres"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

type Suite struct {
	postgresContainer *postgres.PostgresContainer
	kafkaContainer    *kafka.KafkaContainer
	testDelay         time.Duration
	token             string
}

func New(t *testing.T, testDelay time.Duration, token string) *Suite {
	pc := initPostgres(t, context.Background(), "company-db", "user", "passwd")
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(pc); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	})

	kc := initKafka(t, context.Background())
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(kc); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	})

	return &Suite{
		postgresContainer: pc,
		kafkaContainer:    kc,
		testDelay:         testDelay,
		token:             token,
	}
}

func initPostgres(t *testing.T, ctx context.Context, dbName, dbUser, dbPassword string) *postgres.PostgresContainer {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts("../../scripts/init.sql"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	return postgresContainer
}

func initKafka(t *testing.T, ctx context.Context) *kafka.KafkaContainer {
	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
	)
	require.NoError(t, err)

	return kafkaContainer
}

func (s *Suite) testClientPost(t *testing.T, token, url string, data []byte) ([]byte, int) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return body, resp.StatusCode
}

func (s *Suite) testClientGet(t *testing.T, token, url string) ([]byte, int) {
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return body, resp.StatusCode
}

func (s *Suite) testClientDelete(t *testing.T, token, url string) ([]byte, int) {
	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return body, resp.StatusCode
}

func (s *Suite) testClientPatch(t *testing.T, token, url string, data []byte) ([]byte, int) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return body, resp.StatusCode
}

func (s *Suite) startTests(t *testing.T, pg *pkgPg.Postgres, log *logger.Logger) {
	time.Sleep(s.testDelay)

	t.Parallel()
	t.Run("Test CompanyCreate", func(t *testing.T) {
		s.testCreateHttpCases(t, pg, log)
	})

	t.Run("Test CompanyGet", func(t *testing.T) {
		s.testGetHttpCases(t, pg, log)
	})

	t.Run("Test CompanyDelete", func(t *testing.T) {
		s.testDeleteHttpCases(t, pg, log)
	})

	t.Run("Test CompanyPatch", func(t *testing.T) {
		s.testPatchHttpCases(t, pg, log)
	})
}
