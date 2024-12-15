package services

import (
	"company-crud/internal/domain"
	"company-crud/pkg/logger"
	"company-crud/pkg/producer"
	"fmt"
	"github.com/google/uuid"
)

const errorSection = "companyService"
const (
	create  = "create"
	deleteM = "delete"
	get     = "get"
	patch   = "patch"
)

type Company struct {
	producer  *producer.KafkaProducer
	companyDB domain.CompanyDB
	logger    *logger.Logger
}

func New(log *logger.Logger, prod *producer.KafkaProducer, compDB domain.CompanyDB) *Company {
	return &Company{
		producer:  prod,
		companyDB: compDB,
		logger:    log,
	}
}

func (c *Company) Create(company domain.Company) (uuid.UUID, error) {
	id, err := c.companyDB.Insert(company)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = c.producer.ProduceEvent([]byte(fmt.Sprintf("New company created, with ID: %s", id.String())))
	if err != nil {
		return uuid.UUID{}, err
	}

	c.logger.Named(fmt.Sprintf("%s:%s", errorSection, create)).Info("Company entry and event created")

	return id, nil
}

func (c *Company) Delete(companyName string) error {
	err := c.companyDB.DeleteByName(companyName)
	if err != nil {
		return err
	}

	err = c.producer.ProduceEvent([]byte(fmt.Sprintf("New company created, with name: %s", companyName)))
	if err != nil {
		return err
	}

	c.logger.Named(fmt.Sprintf("%s:%s", errorSection, deleteM)).Info("Company entry deleted")

	return nil
}

func (c *Company) Get(companyName string) (domain.Company, error) {
	company, err := c.companyDB.GetByName(companyName)
	if err != nil {
		return domain.Company{}, err
	}

	c.logger.Named(fmt.Sprintf("%s:%s", errorSection, get)).Info("Company info retrieved")

	return company, nil
}

func (c *Company) Patch(company domain.Company, currentName string) error {
	err := c.companyDB.PatchByName(company, currentName)
	if err != nil {
		return err
	}

	err = c.producer.ProduceEvent([]byte(fmt.Sprintf("Company entry patched, with ID: %s", company.ID.String())))
	if err != nil {
		return err
	}

	c.logger.Named(fmt.Sprintf("%s:%s", errorSection, patch)).Info("Company info retrieved")

	return nil
}
