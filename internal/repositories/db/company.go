package db

import (
	"company-crud/internal/domain"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"fmt"
	"github.com/google/uuid"
)

const errorSection = "userDB"
const (
	createMethod  = "create"
	getByIDMethod = "get"
	patchMethod   = "patch"
	deleteMethod  = "delete"
)

type Company struct {
	db     *postres.Postgres
	logger *logger.Logger
}

func New(db *postres.Postgres, log *logger.Logger) *Company {
	return &Company{
		db:     db,
		logger: log,
	}
}

func (u *Company) Insert(company domain.Company) (uuid.UUID, error) {
	userModel := modelConverter(company)

	err := u.db.QueryRow(
		`INSERT INTO companies (name, description, employees_number, is_registered, type, updated_at)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id`,
		userModel.Name,
		userModel.Description,
		userModel.EmployeesNumber,
		userModel.IsRegistered,
		userModel.Type,
		userModel.UpdatedAt,
	).Scan(&userModel.ID)

	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, createMethod)).Error(err.Error())
		return uuid.Nil, err
	}
	return userModel.ID, nil
}
