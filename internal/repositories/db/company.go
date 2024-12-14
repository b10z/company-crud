package db

import (
	"company-crud/internal/domain"
	"company-crud/pkg/logger"
	"company-crud/pkg/postres"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const errorSection = "userDB"
const (
	create       = "create"
	getByName    = "getByName"
	patchByName  = "patchByName"
	deleteByName = "deleteByName"
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
	companyModel := modelConverter(company)

	err := u.db.QueryRow(
		`INSERT INTO xm_assessment.companies (name, description, employees_number, is_registered, type, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 RETURNING id`,
		companyModel.Name,
		companyModel.Description,
		companyModel.EmployeesNumber,
		companyModel.IsRegistered,
		companyModel.Type,
		companyModel.UpdatedAt,
	).Scan(&companyModel.ID)

	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, create)).Error(err.Error())

		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			if pgErr.Code == "23505" {
				return uuid.UUID{}, postres.DuplicateKey
			}
		}

		return uuid.Nil, err
	}
	return companyModel.ID, nil
}

func (u *Company) GetByName(name string) (domain.Company, error) {
	companyModel := model{
		Name: name,
	}

	query := `SELECT * FROM xm_assessment.companies WHERE name = $1`

	var tempType string
	err := u.db.QueryRow(query, companyModel.Name).Scan(&companyModel.ID, &companyModel.Name, &companyModel.Description, &companyModel.EmployeesNumber, &companyModel.IsRegistered, &tempType, &companyModel.CreatedAt, &companyModel.UpdatedAt)
	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, getByName)).Error(err.Error())
		return domain.Company{}, err
	}

	companyType, err := domain.GetCompTypeFromString(tempType)
	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, getByName)).Error(err.Error())
		return domain.Company{}, err
	}

	return domain.Company{
		ID:              companyModel.ID,
		Name:            companyModel.Name,
		Description:     &companyModel.Description,
		EmployeesNumber: &companyModel.EmployeesNumber,
		IsRegistered:    &companyModel.IsRegistered,
		Type:            &companyType,
		UpdatedAt:       companyModel.UpdatedAt,
		CreatedAt:       companyModel.CreatedAt,
	}, nil
}

func (u *Company) DeleteByName(name string) error {
	query := `DELETE FROM xm_assessment.companies WHERE name=$1`

	res, err := u.db.Exec(query, name)
	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, deleteByName)).Error(err.Error())
		return err
	}

	rowsNumber, err := res.RowsAffected()
	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, deleteByName)).Error(err.Error())
		return err
	}

	if rowsNumber == 0 {
		return postres.NoRowsErr
	}

	return nil
}

func (u *Company) PatchByName(company domain.Company, currentName string) error {
	affectedFields, err := patchQueryBuilder(company)
	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, patchByName)).Error(err.Error())
		return err
	}

	companyModel := modelConverter(company)
	companyModel.CurrentName = currentName

	query := fmt.Sprintf(`UPDATE xm_assessment.companies SET %s WHERE name=:current_name`, affectedFields)

	result, err := u.db.NamedExec(query, companyModel)
	if err != nil {
		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, patchByName)).Error(err.Error())
		return err
	}

	rowsNumber, err := result.RowsAffected()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return postres.NoRowsErr
		}

		u.logger.Named(fmt.Sprintf("%s:%s", errorSection, patchByName)).Error(err.Error())
		return err
	}

	if rowsNumber == 0 {
		return postres.NoRowsErr
	}

	return nil
}

func patchQueryBuilder(company domain.Company) (string, error) {
	var query string
	if company.Name != "" {
		query += ` name=:name,`
	}

	if company.Description != nil {
		query += ` description=:description,`
	}
	if company.IsRegistered != nil {
		query += ` is_registered=:is_registered,`
	}

	if company.EmployeesNumber != nil {
		query += ` employees_number=:employees_number,`
	}

	if company.Type != nil {
		query += ` type=:type,`
	}

	if query == "" {
		return query, postres.InvalidArgumentsForBuildingquery
	}
	query += ` updated_at=:updated_at`

	return query, nil
}
