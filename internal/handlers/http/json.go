package http

type Create struct {
	Name              string `json:"username" validate:"required"`
	Description       string `json:"description"`
	AmountOfEmployees int    `json:"amount_of_employees" validate:"required"`
	IsRegistered      bool   `json:"registered" validate:"required"`
	Type              string `json:"type" validate:"required"`
}
