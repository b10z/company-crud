package db

import (
	"company-crud/internal/domain"
	"github.com/google/uuid"
	"time"
)

type model struct {
	ID              uuid.UUID          `db:"id"`
	Name            string             `db:"name,omitempty"`
	Description     string             `db:"description,omitempty"`
	EmployeesNumber int                `db:"employees_number,omitempty"`
	IsRegistered    bool               `db:"is_registered,omitempty"`
	Type            domain.CompanyType `db:"type,omitempty"`
	CreatedAt       time.Time          `db:"created_at,omitempty"`
	UpdatedAt       time.Time          `db:"updated_at"`
}

func modelConverter(d domain.Company) model {
	userModel := model{
		ID:              d.ID,
		Name:            d.Name,
		Description:     *d.Description,
		EmployeesNumber: *d.EmployeesNumber,
		IsRegistered:    *d.IsRegistered,
		Type:            *d.Type,
		UpdatedAt:       time.Now(),
	}

	return userModel
}
