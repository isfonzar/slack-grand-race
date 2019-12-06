package message

import "testing"

type (
	LoggerMock struct {
	}
)

func (lm LoggerMock) Debug(args ...interface{})                   {}
func (lm LoggerMock) Info(args ...interface{})                    {}
func (lm LoggerMock) Warn(args ...interface{})                    {}
func (lm LoggerMock) Error(args ...interface{})                   {}
func (lm LoggerMock) Fatal(args ...interface{})                   {}
func (lm LoggerMock) Panic(args ...interface{})                   {}
func (lm LoggerMock) Debugw(template string, args ...interface{}) {}
func (lm LoggerMock) Infow(template string, args ...interface{})  {}
func (lm LoggerMock) Warnw(template string, args ...interface{})  {}
func (lm LoggerMock) Errorw(template string, args ...interface{}) {}

func TestHandler_HandleMessage(t *testing.T) {
	loggerMock := LoggerMock{}
	handler := NewHandler(loggerMock)

	msg := Message{
		channel:   "",
		user:      "",
		content:   "",
		timestamp: "",
	}

	handler.HandleMessage(msg)
}
