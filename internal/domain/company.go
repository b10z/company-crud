package domain

import "github.com/google/uuid"

type Company struct {
	ID              uuid.UUID
	Name            string
	Description     string
	EmployeesNumber int
	IsRegistered    bool
	Type            CompanyType
}

type CompanyType uint8

const (
	Corporations CompanyType = iota
	NonProfit
	Cooperative
	SoleProprietorship
)

func (ct *CompanyType) String() string {
	switch *ct {
	case Corporations:
		return "Corporations"
	case NonProfit:
		return "NonProfit"
	case Cooperative:
		return "Cooperative"
	case SoleProprietorship:
		return "SoleProprietorship"
	default:
		return ""
	}
}

type CompanyDB interface {
	Insert(Company) (uuid.UUID, error)
}
