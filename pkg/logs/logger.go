package logs

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var (
	InitializationError = errors.New("could not initialize logger")
)

// New creates and returns a new Logger instance.
func New(isDebug bool) (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	if !isDebug {
		config = zap.NewProductionConfig()
	}

	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("%w : %v", InitializationError, err)
	}

	return logger.Sugar(), nil
}
