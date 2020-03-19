package domain

import (
	"errors"
	"testing"

	"github.com/nlopes/slack"
)

type (
	UserInfoMock struct {
		u   *slack.User
		err error
	}
)

func (uim *UserInfoMock) GetUserInfo(user string) (*slack.User, error) {
	return uim.u, uim.err
}

func TestNewUserFromSlack(t *testing.T) {
	var tests = []struct {
		id          string
		name        string
		slackErr    error
		expectedErr error
	}{
		{"U5NTYR0EQ", "user_name", nil, nil},
		{"U5NTYR0EQ", "user_name", errors.New("connection issue"), ErrUnableToGetUserInfo},
	}

	for _, test := range tests {
		su := slack.User{
			ID:   test.id,
			Name: test.name,
		}

		uim := &UserInfoMock{
			u:   &su,
			err: test.slackErr,
		}

		u, err := NewUserFromSlack(uim, test.id)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("NewUserFromSlack() returned an unexpected error, got: %v, expected: %v", err, test.expectedErr)
		}
		if err == nil &&
			(u.ID != test.id ||
				u.Name != test.name) {
			t.Errorf("NewUserFromSlack() does not match info got from provider, user: %v, provider: %v", u, test)
		}
	}
}
