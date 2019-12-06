package message

import "github.com/nlopes/slack"

type Message struct {
	Channel   string
	User      string
	Content   string
	Timestamp string
}

func NewMessageFromEvent(slackMessage *slack.MessageEvent) Message {
	return Message{
		Channel:   slackMessage.Channel,
		User:      slackMessage.User,
		Content:   slackMessage.Text,
		Timestamp: slackMessage.Timestamp,
	}
}
