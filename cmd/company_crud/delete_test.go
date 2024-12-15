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

func (s *Suite) testDeleteHttpCases(t *testing.T, pg *pkgPg.Postgres, log *logger.Logger) {
	t.Run("Valid delete - with token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNameDel_1",
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

		_, status = s.testClientDelete(t, s.token, "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusOK, status)

		companyDBData, err = companyDB.GetByName(req.Name)
		require.ErrorIs(t, err, pkgPg.NoRowsErr)

		_, status = s.testClientGet(t, s.token, "http://localhost:8000/companies/"+req.Name)
		require.Equal(t, http.StatusConflict, status)

	})

	t.Run("Valid delete - without token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNameDel_2",
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

		_, status = s.testClientDelete(t, "", "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusForbidden, status)

		companyDBData, err = companyDB.GetByName(req.Name)
		require.NoError(t, err)

		_, status = s.testClientGet(t, s.token, "http://localhost:8000/companies/"+companyDBData.Name)
		require.Equal(t, http.StatusOK, status)
	})

	t.Run("Delete invalid delete", func(t *testing.T) {
		_, status := s.testClientDelete(t, s.token, "http://localhost:8000/companies/"+"randomName")
		require.Equal(t, http.StatusConflict, status)
	})
}
