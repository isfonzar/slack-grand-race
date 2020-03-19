package domain

import (
	"testing"

	"github.com/slack-go/slack"
)

func TestNewMessageFromSlack(t *testing.T) {
	var tests = []struct {
		channel   string
		user      string
		text      string
		timestamp string
	}{
		{"channel", "U5NTYR0EQ", "content", "1584549498.001000"},
	}

	for _, test := range tests {
		ev := slack.MessageEvent{
			Msg: slack.Msg{
				Channel:   test.channel,
				User:      test.user,
				Text:      test.text,
				Timestamp: test.timestamp,
			},
		}

		m := NewMessageFromSlack(&ev)
		if m.User != test.user ||
			m.Channel != test.channel ||
			m.Timestamp != test.timestamp ||
			m.Content != test.text {
			t.Errorf("NewMessageFromSlack() does not match info got from slack, message: %v, provider: %v", m, test)
		}
	}
}
