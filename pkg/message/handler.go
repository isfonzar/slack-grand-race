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

func (h Handler) HandleMessage(message Message) {

}
