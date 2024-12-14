package postres

import (
	"company-crud/pkg/logger"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	InvalidArgumentsForBuildingquery = errors.New("invalid arguments for building a query")
	NoRowsErr                        = errors.New("no rows")
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUsername string
	DBPassword string
	DBName     string
}

type Postgres struct {
	*sqlx.DB
	log *logger.Logger
	cfg Config
}

func New(log *logger.Logger, conf Config) (*Postgres, error) {
	postgresConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost,
		conf.DBPort,
		conf.DBUsername,
		conf.DBPassword,
		conf.DBName,
	)

	db, err := sqlx.Open("postgres", postgresConn)
	if err != nil {
		log.Error("Error opening connection with db", zap.Error(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error("Error while pinging db", zap.Error(err))
		return nil, err
	}

	log.Info("The database is connected")

	return &Postgres{
		db,
		log,
		conf,
	}, err
}
