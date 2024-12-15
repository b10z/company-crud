package main

import (
	jsons "company-crud/internal/handlers/http"
	"company-crud/internal/repositories/db"
	"company-crud/pkg/logger"
	pkgPg "company-crud/pkg/postres"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func (s *Suite) testGetHttpCases(t *testing.T, pg *pkgPg.Postgres, log *logger.Logger) {
	t.Run("Valid get - with token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNameGet_1",
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

		resp, status := s.testClientGet(t, s.token, "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusOK, status)

		getResp := jsons.Get{}
		err = json.Unmarshal(resp, &getResp)
		require.NoError(t, err)

		require.Equal(t, req.Name, getResp.Name)
		require.Equal(t, req.Description, getResp.Description)
		require.Equal(t, req.EmployeesNumber, &getResp.EmployeesNumber)
		require.Equal(t, req.IsRegistered, getResp.IsRegistered)
		require.Equal(t, req.Type, getResp.Type)

	})

	t.Run("Valid get - without token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNameGet_2",
			Description:     "description_2",
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

		_, status = s.testClientGet(t, "", "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusForbidden, status)

	})

	t.Run("Valid get from entry with no description", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNameGet_3",
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

		resp, status := s.testClientGet(t, s.token, "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusOK, status)

		getResp := jsons.Get{}
		err = json.Unmarshal(resp, &getResp)
		require.NoError(t, err)

		require.Equal(t, req.Name, getResp.Name)
		require.Equal(t, req.Description, getResp.Description)
		require.Equal(t, req.EmployeesNumber, &getResp.EmployeesNumber)
		require.Equal(t, req.IsRegistered, getResp.IsRegistered)
		require.Equal(t, req.Type, getResp.Type)

	})

	t.Run("Valid get from entry with 0 employees number", func(t *testing.T) {
		employees := 0
		req := jsons.Create{
			Name:            "testNameGet_4",
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

		resp, status := s.testClientGet(t, s.token, "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusOK, status)

		getResp := jsons.Get{}
		err = json.Unmarshal(resp, &getResp)
		require.NoError(t, err)

		require.Equal(t, req.Name, getResp.Name)
		require.Equal(t, req.Description, getResp.Description)
		require.Equal(t, req.EmployeesNumber, &getResp.EmployeesNumber)
		require.Equal(t, req.IsRegistered, getResp.IsRegistered)
		require.Equal(t, req.Type, getResp.Type)
	})

	t.Run("Get invalid name", func(t *testing.T) {
		_, status := s.testClientGet(t, s.token, "http://localhost:8000/companies/"+"randomName")
		require.Equal(t, http.StatusConflict, status)
	})
}
