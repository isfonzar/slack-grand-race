package message

import (
	"github.com/nlopes/slack"
	"testing"
)

func TestMessage(t *testing.T) {
	slackMsg := slack.Msg{
		User:      "I4GH53KAS9",
		Channel:   "G72NC8JS98",
		Timestamp: "1575637406.001000",
		Text:      "test message",
	}

	slackMessageEvent := slack.MessageEvent{
		Msg:        slackMsg,
		SubMessage: nil,
	}

	message := NewMessageFromEvent(&slackMessageEvent)

	if message.GetUser() != slackMsg.User {
		t.Error("User is not set correctly")
	}
	if message.GetChannel() != slackMsg.Channel {
		t.Error("Channel is not set correctly")
	}
	if message.getContent() != slackMsg.Text {
		t.Error("Message text content is not set correctly")
	}
	if message.GetTimestamp() != slackMsg.Timestamp {
		t.Error("Timestamp is not set correctly")
	}
}
