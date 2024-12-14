package http

import (
	"company-crud/internal/domain"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"company-crud/pkg/validator"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const errorSection = "companyHandler"
const (
	create  = "create"
	deleteM = "deleteM"
	patch   = "patch"
	get     = "get"
)

type Company struct {
	logger         *logger.Logger
	validator      *validator.Validator
	companyService domain.CompanyService
	tokenSignature string
}

func New(log *logger.Logger, cs domain.CompanyService, tokenSig string) *Company {
	return &Company{
		logger:         log,
		validator:      validator.New(),
		companyService: cs,
		tokenSignature: tokenSig,
	}
}

func (c *Company) AddRoute(r *mux.Router) {
	companiesRoutes := r.PathPrefix("/companies").Subrouter()
	companiesRoutes.Use(validateToken(c.tokenSignature))
	companiesRoutes.HandleFunc("", c.create).Methods(http.MethodPost)
	companiesRoutes.HandleFunc("/{company_name}", c.get).Methods(http.MethodGet)
	companiesRoutes.HandleFunc("/{company_name}", c.delete).Methods(http.MethodDelete)
	companiesRoutes.HandleFunc("/{company_name}", c.patch).Methods(http.MethodPatch)
}

// @Summary      Create new company
// @Tags         company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        createCompany	body	Create  true  "createCompany"
// @Success      200
// @Failure      400
// @Failure      406
// @Failure      409	{object}  Error true
// @Router       /companies [post]
func (c *Company) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqData := Create{}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, create)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(reqData); err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, create)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	companyType, err := domain.GetCompTypeFromString(reqData.Type)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, create)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = c.companyService.Create(domain.Company{
		Name:            reqData.Name,
		Description:     &reqData.Description,
		EmployeesNumber: reqData.EmployeesNumber,
		IsRegistered:    &reqData.IsRegistered,
		Type:            &companyType,
	})
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, create)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)

		if errors.As(err, &postres.DuplicateKey) {
			w.WriteHeader(http.StatusConflict)
			resp, _ := json.Marshal(Error{
				Message: "duplicate name",
			})
			w.Write(resp)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Get company
// @Tags         company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        company_name	path	string true "company_name"
// @Success      200
// @Failure      400
// @Failure      406
// @Failure      409	{object}  Error true
// @Failure      500
// @Router       /companies/{company_name} [get]
func (c *Company) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nameParam := mux.Vars(r)["company_name"]

	if nameParam == "" {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, get)).Debug("company_name param can't be empty")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	result, err := c.companyService.Get(nameParam)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, get)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &postres.NoRowsErr) {
			w.WriteHeader(http.StatusConflict)
			resp, _ := json.Marshal(Error{
				Message: "no results",
			})
			w.Write(resp)
		}

		return
	}

	response, err := json.Marshal(Get{
		ID:              result.ID,
		Name:            result.Name,
		Description:     *result.Description,
		EmployeesNumber: *result.EmployeesNumber,
		IsRegistered:    *result.IsRegistered,
		Type:            result.Type.String(),
		UpdatedAt:       result.UpdatedAt,
		CreatedAt:       result.CreatedAt,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// @Summary      Delete company
// @Tags         company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        company_name	path	string true "company_name"
// @Success      200
// @Failure      400
// @Failure      406
// @Failure      409	{object}  Error true
// @Failure      500
// @Router       /companies/{company_name} [delete]
func (c *Company) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nameParam := mux.Vars(r)["company_name"]

	if nameParam == "" {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, deleteM)).Debug("company_name param can't be empty")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err := c.companyService.Delete(nameParam)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, deleteM)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &postres.NoRowsErr) {
			w.WriteHeader(http.StatusConflict)
			resp, _ := json.Marshal(Error{
				Message: "no entries affected",
			})
			w.Write(resp)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Patch company
// @Tags         company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        company_name	path	string true "company_name"
// @Param        patchCompany	body	Patch  true  "patchCompany"
// @Success      200
// @Failure      400
// @Failure      406
// @Failure      409	{object}  Error true
// @Failure      500
// @Router       /companies/{company_name} [patch]
func (c *Company) patch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nameParam := mux.Vars(r)["company_name"]

	if nameParam == "" {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, patch)).Debug("company_name param can't be empty")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	reqData := Patch{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, patch)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(reqData); err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, patch)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	companyPatch := domain.Company{
		Name:            reqData.Name,
		Description:     reqData.Description,
		EmployeesNumber: reqData.EmployeesNumber,
		IsRegistered:    reqData.IsRegistered,
	}
	if reqData.Type != nil {
		companyType, err := domain.GetCompTypeFromString(*reqData.Type)
		if err != nil {
			c.logger.Named(fmt.Sprintf("%s:%s", errorSection, patch)).Debug(err.Error())
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		companyPatch.Type = &companyType
	} else {
		companyPatch.Type = nil
	}

	err = c.companyService.Patch(companyPatch, nameParam)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, patch)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &postres.NoRowsErr) {
			w.WriteHeader(http.StatusConflict)
			resp, _ := json.Marshal(Error{
				Message: "no entries affected",
			})
			w.Write(resp)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
