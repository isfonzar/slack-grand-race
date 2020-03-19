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
