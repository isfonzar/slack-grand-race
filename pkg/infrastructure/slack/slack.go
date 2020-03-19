package slack

import (
	"github.com/isfonzar/slack-grand-race/pkg/domain"
	"github.com/slack-go/slack"
)

type (
	Handler struct {
		sl Slacker
	}

	Slacker interface {
		AddReaction(name string, item slack.ItemRef) error
		SendMessage(msg *slack.OutgoingMessage)
		NewOutgoingMessage(text string, channelID string, options ...slack.RTMsgOption) *slack.OutgoingMessage
	}
)

func NewHandler(sl Slacker) *Handler {
	return &Handler{
		sl: sl,
	}
}

func (h *Handler) AddReaction(msg *domain.Message, reaction domain.Reaction) error {
	return h.sl.AddReaction(string(reaction), slack.NewRefToMessage(msg.Channel, msg.Timestamp))
}

func (h *Handler) SendMessage(text string, msg *domain.Message) {
	h.sl.SendMessage(h.sl.NewOutgoingMessage(text, msg.Channel))
}
