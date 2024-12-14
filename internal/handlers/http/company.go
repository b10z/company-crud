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
	createMethod = "createCompany"
)

type Company struct {
	logger         *logger.Logger
	validator      *validator.Validator
	companyService domain.CompanyService
}

func New(log *logger.Logger, cs domain.CompanyService) *Company {
	return &Company{
		logger:         log,
		companyService: cs,
		validator:      validator.New(),
	}
}

func (c *Company) AddRoute(r *mux.Router) {
	companiesRoutes := r.PathPrefix("/companies").Subrouter()
	companiesRoutes.HandleFunc("", c.create).Methods(http.MethodPost)
}

// @Summary      Create new company
// @Tags         company
// @Accept       json
// @Produce      json
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
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(reqData); err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	companyType, err := domain.GetCompTypeFromString(reqData.Type)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = c.companyService.Create(domain.Company{
		Name:            reqData.Name,
		Description:     &reqData.Description,
		EmployeesNumber: &reqData.EmployeesNumber,
		IsRegistered:    &reqData.IsRegistered,
		Type:            &companyType,
	})
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)

		if errors.As(err, &postres.NoRowsErr) {
			w.WriteHeader(http.StatusConflict)
			resp, _ := json.Marshal(Error{
				Message: "duplicate name ",
			})
			w.Write(resp)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Get new company
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_name	path	string true "company_name"
// @Success      200
// @Failure      400
// @Failure      406
// @Failure      409	{object}  Error true
// @Router       /companies/{company_name} [get]
func (c *Company) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqData := Create{}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(reqData); err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	companyType, err := domain.GetCompTypeFromString(reqData.Type)
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = c.companyService.Create(domain.Company{
		Name:            reqData.Name,
		Description:     &reqData.Description,
		EmployeesNumber: &reqData.EmployeesNumber,
		IsRegistered:    &reqData.IsRegistered,
		Type:            &companyType,
	})
	if err != nil {
		c.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Debug(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)

		if errors.As(err, &postres.NoRowsErr) {
			w.WriteHeader(http.StatusConflict)
			resp, _ := json.Marshal(Error{
				Message: "duplicate name ",
			})
			w.Write(resp)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
