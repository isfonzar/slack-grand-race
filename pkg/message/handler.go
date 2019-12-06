package message

type (
	logger interface {
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

	Handler struct {
		logger
	}
)

// NewHandler returns a new message handler
func NewHandler(logger logger) Handler {
	return Handler{logger}
}

// HandleMessage handles a message
func (h Handler) HandleMessage(msg Message) error {
	fields := []interface{}{"message", msg}
	h.logger.Debugw("Message received", fields...)

	// Discard messages from bots (including itself)

	// Is the message an action call?

	// Does the message contain a banned word?

	// Does the message contain a praised word?

	// Should the message receive a reward?

	return nil
}
