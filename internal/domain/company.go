package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Company struct {
	ID              uuid.UUID
	Name            string
	Description     *string
	EmployeesNumber *int
	IsRegistered    *bool
	Type            *CompanyType
	UpdatedAt       time.Time
	CreatedAt       time.Time
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
		return "Sole Proprietorship"
	default:
		return ""
	}
}

func GetCompTypeFromString(compType string) (CompanyType, error) {
	switch compType {
	case "Corporations":
		return Corporations, nil
	case "NonProfit":
		return NonProfit, nil
	case "Cooperative":
		return Cooperative, nil
	case "Sole Proprietorship":
		return SoleProprietorship, nil
	default:
		return 0, fmt.Errorf("invalid companyType")
	}
}

type CompanyDB interface {
	Insert(Company) (uuid.UUID, error)
	DeleteByName(string) error
	PatchByName(Company, string) error
	GetByName(string) (Company, error)
}

type CompanyService interface {
	Create(Company) (uuid.UUID, error)
	Delete(string) error
	Patch(Company, string) error
	Get(string) (Company, error)
}
