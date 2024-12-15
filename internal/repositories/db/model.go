package db

import (
	"company-crud/internal/domain"
	"github.com/google/uuid"
	"time"
)

type model struct {
	ID              uuid.UUID `db:"id"`
	CurrentName     string    `db:"current_name,omitempty"`
	Name            string    `db:"name,omitempty"`
	Description     string    `db:"description,omitempty"`
	EmployeesNumber int       `db:"employees_number,omitempty"`
	IsRegistered    bool      `db:"is_registered,omitempty"`
	Type            string    `db:"type,omitempty"`
	CreatedAt       time.Time `db:"created_at,omitempty"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func modelConverter(d domain.Company) model {
	companyModel := model{
		ID:        d.ID,
		Name:      d.Name,
		UpdatedAt: time.Now(),
	}

	if d.Type != nil {
		companyModel.Type = d.Type.String()
	}

	if d.Description != nil {
		companyModel.Description = *d.Description
	}

	if d.EmployeesNumber != nil {
		companyModel.EmployeesNumber = *d.EmployeesNumber
	}

	if d.IsRegistered != nil {
		companyModel.IsRegistered = *d.IsRegistered
	}

	return companyModel
}
