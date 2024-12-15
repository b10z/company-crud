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

func (s *Suite) testPatchHttpCases(t *testing.T, pg *pkgPg.Postgres, log *logger.Logger) {
	t.Run("Valid patch - with token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNamePatch_1",
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
		_, err = companyDB.GetByName(req.Name)
		require.NoError(t, err)

		employees = 1
		isRegistered := false
		companyType := "Sole Proprietorship"
		desc := "test_patched_desc_1"
		postReq := jsons.Patch{
			Name:            "test_patched1",
			Description:     &desc,
			EmployeesNumber: &employees,
			IsRegistered:    &isRegistered,
			Type:            &companyType,
		}

		jsonData, err = json.Marshal(postReq)
		require.NoError(t, err)

		_, status = s.testClientPatch(t, s.token, "http://localhost:8000/companies/"+req.Name, jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDBData, err := companyDB.GetByName(postReq.Name)
		require.NoError(t, err)

		require.Equal(t, postReq.Name, companyDBData.Name)
		require.Equal(t, postReq.Description, companyDBData.Description)
		require.Equal(t, postReq.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, postReq.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, *postReq.Type, companyDBData.Type.String())
	})

	t.Run("Valid patch - without token", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNamePatch_2",
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
		_, err = companyDB.GetByName(req.Name)
		require.NoError(t, err)

		employees = 1
		isRegistered := false
		companyType := "Sole Proprietorship"
		desc := "test_patched_desc_1"
		postReq := jsons.Patch{
			Name:            "test_patched2",
			Description:     &desc,
			EmployeesNumber: &employees,
			IsRegistered:    &isRegistered,
			Type:            &companyType,
		}

		jsonData, err = json.Marshal(postReq)
		require.NoError(t, err)

		_, status = s.testClientPatch(t, "", "http://localhost:8000/companies/"+req.Name, jsonData)
		require.Equal(t, http.StatusForbidden, status)

		_, err = companyDB.GetByName(postReq.Name)
		require.ErrorIs(t, err, pkgPg.NoRowsErr)
	})

	t.Run("Patch name", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNamePatch_3",
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
		_, err = companyDB.GetByName(req.Name)
		require.NoError(t, err)

		postReq := jsons.Patch{
			Name: "test_patched3",
		}

		jsonData, err = json.Marshal(postReq)
		require.NoError(t, err)

		_, status = s.testClientPatch(t, s.token, "http://localhost:8000/companies/"+req.Name, jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDBData, err := companyDB.GetByName(postReq.Name)
		require.NoError(t, err)

		require.Equal(t, postReq.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, req.Type, companyDBData.Type.String())
	})

	t.Run("Patch type", func(t *testing.T) {
		employees := 2
		req := jsons.Create{
			Name:            "testNamePatch_4",
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
		_, err = companyDB.GetByName(req.Name)
		require.NoError(t, err)

		newType := "Sole Proprietorship"
		postReq := jsons.Patch{
			Type: &newType,
		}

		jsonData, err = json.Marshal(postReq)
		require.NoError(t, err)

		_, status = s.testClientPatch(t, s.token, "http://localhost:8000/companies/"+req.Name, jsonData)
		require.Equal(t, http.StatusOK, status)

		companyDBData, err := companyDB.GetByName(req.Name)
		require.NoError(t, err)

		require.Equal(t, req.Name, companyDBData.Name)
		require.Equal(t, &req.Description, companyDBData.Description)
		require.Equal(t, req.EmployeesNumber, companyDBData.EmployeesNumber)
		require.Equal(t, &req.IsRegistered, companyDBData.IsRegistered)
		require.Equal(t, *postReq.Type, companyDBData.Type.String())
	})

	t.Run("Patch with invalid json", func(t *testing.T) {
		_, status := s.testClientPatch(t, s.token, "http://localhost:8000/companies/"+"randomName", []byte{})
		require.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("Patch with invalid name", func(t *testing.T) {
		postReq := jsons.Patch{
			Name: "test_patched3",
		}
		jsonData, err := json.Marshal(postReq)
		require.NoError(t, err)

		_, status := s.testClientPatch(t, s.token, "http://localhost:8000/companies/"+"randomName", jsonData)
		require.Equal(t, http.StatusConflict, status)
	})
}
