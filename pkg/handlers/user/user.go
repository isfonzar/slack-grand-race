package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
	"github.com/slack-go/slack"
)

type (
	Handler struct {
		info    Info
		storage Storage
	}

	Info interface {
		GetUserInfo(user string) (*slack.User, error)
	}

	Storage interface {
		Get(id string) (*domain.User, error)
		Create(id, name string) error
	}
)

var (
	ErrQueryDatabase       = errors.New("could not fetch user from database")
	ErrCreateUser          = errors.New("could not create user in database")
	ErrUnableToGetUserInfo = errors.New("NewUserFromSlack could not get user info")
)

func NewHandler(info Info, storage Storage) *Handler {
	return &Handler{
		info:    info,
		storage: storage,
	}
}

func (h *Handler) GetUser(id string) (*domain.User, error) {
	// Get user from database
	u, err := h.storage.Get(id)
	if err != nil {
		return u, fmt.Errorf("%w : %v", ErrQueryDatabase, err)
	}

	// User does not exist, create it.
	if u == nil {
		su, err := h.getFromSlack(id)
		if err != nil {
			return u, fmt.Errorf("%w: %v", ErrUnableToGetUserInfo, err)
		}

		if err := h.storage.Create(su.Id, su.Name); err != nil {
			return u, fmt.Errorf("%w: %v", ErrCreateUser, err)
		}

		// Wait so it does not break the select
		time.Sleep(3 * time.Second)

		u, err := h.storage.Get(su.Id)
		if err != nil || u == nil {
			return u, fmt.Errorf("%w : %v", ErrQueryDatabase, err)
		}
	}

	return u, nil
}

func (h *Handler) getFromSlack(id string) (*domain.User, error) {
	var user domain.User

	su, err := h.info.GetUserInfo(id)
	if err != nil {
		return &user, err
	}

	return &domain.User{
		Id:   su.ID,
		Name: su.Name,
	}, nil
}
