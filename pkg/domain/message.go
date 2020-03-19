package domain

import (
	"github.com/slack-go/slack"
)

type (
	Message struct {
		Channel   string
		User      string
		Content   string
		Timestamp string
	}
)

func NewMessageFromSlack(ev *slack.MessageEvent) *Message {
	return &Message{
		Channel:   ev.Channel,
		User:      ev.User,
		Content:   ev.Text,
		Timestamp: ev.Timestamp,
	}
}
