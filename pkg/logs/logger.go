package logs

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	//Logger is the interface used by the logger client.
	Logger interface {
		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})
		Debugw(template string, args ...interface{})
		Infow(template string, args ...interface{})
		Warnw(template string, args ...interface{})
		Errorw(template string, args ...interface{})
	}
)

//New creates and returns a new Logger instance.
func New(isDebug bool) (Logger, error) {
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
