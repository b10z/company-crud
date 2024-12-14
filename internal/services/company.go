package services

import (
	"company-crud/internal/domain"
	"company-crud/pkg/logger"
	"company-crud/pkg/producer"
	"fmt"
	"github.com/google/uuid"
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

	err = c.producer.ProduceEvent([]byte(fmt.Sprintf("New id")))

}
