package message

import "github.com/nlopes/slack"

type Message struct {
	channel   string
	user      string
	content   string
	timestamp string
}

func NewMessageFromEvent(slackMessage *slack.MessageEvent) Message {
	return Message{
		channel:   slackMessage.Channel,
		user:      slackMessage.User,
		content:   slackMessage.Text,
		timestamp: slackMessage.Timestamp,
	}
}

func (m Message) GetChannel() string {
	return m.channel
}

func (m Message) GetUser() string {
	return m.user
}

func (m Message) getContent() string {
	return m.content
}

func (m Message) GetTimestamp() string {
	return m.timestamp
}
