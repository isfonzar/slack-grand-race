package message

import (
	"fmt"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
)

type (
	Handler struct {
	}
)

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Process(msg *domain.Message, user *domain.User) error {
	fmt.Println(msg)
	fmt.Println(user)

	return nil
}
