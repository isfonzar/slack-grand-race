package domain

import (
	"errors"
	"fmt"

	"github.com/nlopes/slack"
)

type (
	User struct {
		ID   string
		Name string
	}

	UserInfo interface {
		GetUserInfo(user string) (*slack.User, error)
	}
)

var (
	ErrUnableToGetUserInfo = errors.New("NewUserFromSlack could not get user info")
)

func NewUserFromSlack(info UserInfo, id string) (*User, error) {
	var user User

	su, err := info.GetUserInfo(id)
	if err != nil {
		return &user, fmt.Errorf("%w: %v", ErrUnableToGetUserInfo, err)
	}

	return &User{
		ID:   su.ID,
		Name: su.Name,
	}, nil
}
