package logs

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
		return nil, errors.Wrap(err, "error initializing zap logger")
	}

	return logger.Sugar(), nil
}
