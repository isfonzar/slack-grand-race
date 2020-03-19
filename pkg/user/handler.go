package user

import (
	"github.com/isfonzar/slack-grand-race/pkg/domain"
)

type (
	Handler struct {
	}

	Slack interface {
		GetUserInfo(user string) (*domain.User, error)
	}
)
