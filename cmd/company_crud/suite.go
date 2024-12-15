package main

import (
	"bytes"
	"company-crud/internal/domain"
	jsons "company-crud/internal/handlers/http"
	"company-crud/internal/repositories/db"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"context"
	"encoding/json"
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

func testClientPost(t *testing.T, token, url string, data []byte) ([]byte, int) {
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

func testClientGet(t *testing.T, token, url string, data []byte) ([]byte, int) {
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

func testClientDelete(t *testing.T, token, url string, data []byte) ([]byte, int) {
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

func testClientPatch(t *testing.T, token, url string, data []byte) ([]byte, int) {
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

func (s *Suite) testCreateHttpCases(t *testing.T, pg *postres.Postgres, log *logger.Logger) {
	time.Sleep(s.testDelay)

	t.Run("Valid creation - with token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_1",
			Description:     "description_1",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "NonProfit",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Valid creation - without token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_2",
			Description:     "description_2",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "NonProfit",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, "", "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusForbidden, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, postres.NoRowsErr)
		require.Equal(t, companyDBData, domain.Company{})
	})

	t.Run("Invalid company type for creation", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_3",
			Description:     "description_3",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "invalidType",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusNotAcceptable, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, postres.NoRowsErr)
		require.Equal(t, companyDBData, domain.Company{})
	})

	t.Run("Invalid Creation with missing name", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Description:     "description_4",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "invalidType",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusNotAcceptable, status)
	})

	t.Run("Invalid Creation with missing description", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_5",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "NonProfit",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Invalid Creation with missing employeesNumber", func(t *testing.T) {
		req := jsons.Create{
			Name:         "testName_6",
			IsRegistered: true,
			Type:         "NonProfit",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusNotAcceptable, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Invalid Creation with missing isRegistered", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_7",
			Description:     "description_7",
			EmployeesNumber: &employees,
			Type:            "NonProfit",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, "", "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusForbidden, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, postres.NoRowsErr)
		require.Equal(t, companyDBData, domain.Company{})
	})

	t.Run("Invalid Creation with missing type", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_8",
			Description:     "description_8",
			EmployeesNumber: &employees,
			IsRegistered:    true,
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, "", "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusForbidden, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, postres.NoRowsErr)
		require.Equal(t, companyDBData, domain.Company{})
	})

	t.Run("Valid creation with Corporations type", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_9",
			Description:     "description_9",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "Corporations",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Valid creation with Cooperative type", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_10",
			Description:     "description_10",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "Cooperative",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Valid creation with Sole Proprietorship type", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_11",
			Description:     "description_11",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "Sole Proprietorship",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Valid creation with employeesNumber 0", func(t *testing.T) {
		employees := 0
		req := jsons.Create{
			Name:            "testName_12",
			Description:     "description_12",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "Sole Proprietorship",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)
		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})
}
