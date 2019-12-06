package message

type (
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

	Handler struct {
		Logger
	}
)

func NewHandler(logger Logger) Handler {
	return Handler{logger}
}

func (h Handler) HandleMessage(msg Message) error {
	fields := []interface{}{"message", msg}
	h.Logger.Debugw("Message received", fields...)

	// Discard messages from the bot itself

	// Discard messages from other bots

	// Is the message an action call?

	// Does the message contain a banned word?

	// Does the message contain a praised word?

	// Should the message receive a reward?

	return nil
}
