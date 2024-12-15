package http

import (
	"github.com/google/uuid"
	"time"
)

type Create struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	EmployeesNumber *int   `json:"amount_of_employees" validate:"required"`
	IsRegistered    bool   `json:"registered" validate:"required"`
	Type            string `json:"type" validate:"required"`
}

type Get struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	EmployeesNumber int       `json:"amount_of_employees"`
	IsRegistered    bool      `json:"registered"`
	Type            string    `json:"type"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type Patch struct {
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	EmployeesNumber *int    `json:"amount_of_employees"`
	IsRegistered    *bool   `json:"registered"`
	Type            *string `json:"type"`
}

type Error struct {
	Message string `json:"error_message"`
}
