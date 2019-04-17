package logs

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	//Logger is the interface used by the logger client.
	Logger interface {
		Info(args ...interface{})
		Infow(msg string, keysAndValues ...interface{})
		Warn(args ...interface{})
		Warnw(msg string, keysAndValues ...interface{})
	}
)

//New creates and returns a new Logger instance.
func New() (Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error initializing zap logger")
	}

	return logger.Sugar(), nil
}
