package main

import (
	"company-crud/internal/domain"
	jsons "company-crud/internal/handlers/http"
	"company-crud/internal/repositories/db"
	"company-crud/pkg/logger"
	pkgPg "company-crud/pkg/postres"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func (s *Suite) testCreateHttpCases(t *testing.T, pg *pkgPg.Postgres, log *logger.Logger) {
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
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

	t.Run("Valid creation but duplicate - with token", func(t *testing.T) {
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusConflict, status)
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

		_, status := s.testClientPost(t, "", "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusForbidden, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, pkgPg.NoRowsErr)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusNotAcceptable, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, pkgPg.NoRowsErr)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusNotAcceptable, status)
	})

	t.Run("Valid Creation with missing description", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testName_5",
			EmployeesNumber: &employees,
			IsRegistered:    true,
			Type:            "NonProfit",
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusNotAcceptable, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, pkgPg.NoRowsErr, err)
		require.Equal(t, domain.Company{}, companyDBData)
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

		_, status := s.testClientPost(t, "", "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusForbidden, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, pkgPg.NoRowsErr)
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

		_, status := s.testClientPost(t, "", "http://localhost:8000/companies", jsonData)
		require.Equal(t, http.StatusForbidden, status)

		companyDB := db.New(pg, log)
		companyDBData, err := companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, pkgPg.NoRowsErr)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
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

		_, status := s.testClientPost(t, s.token, "http://localhost:8000/companies", jsonData)
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
