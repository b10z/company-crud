package logger

import (
	"errors"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New(env string) (*Logger, error) {
	logger := zap.NewNop()

	var err error
	switch env {
	case "test":
		return &Logger{logger}, nil
	case "stag":
		fallthrough
	case "dev":
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

	case "prod":
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("invalid env")
	}

	return &Logger{logger}, nil
}
