package postres

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DuplicateKey                     = errors.New("duplicateKey")
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
	cfg Config
}

func New(conf Config) (*Postgres, error) {
	postgresConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost,
		conf.DBPort,
		conf.DBUsername,
		conf.DBPassword,
		conf.DBName,
	)

	db, err := sqlx.Open("postgres", postgresConn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection with db %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db %w", err)
	}

	return &Postgres{
		db,
		conf,
	}, err
}
