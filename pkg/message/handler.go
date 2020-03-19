package message

import (
	"fmt"

	"github.com/isfonzar/slack-grand-race/internal/repository/postgres"
	"github.com/isfonzar/slack-grand-race/pkg/domain"
	"github.com/pkg/errors"
)

type (
	Handler struct {
		us UserStorage
	}

	UserStorage interface {
		Get(id string) (*postgres.User, error)
		Create(id, name string) error
	}
)

func NewHandler(us UserStorage) *Handler {
	return &Handler{
		us: us,
	}
}

func (h *Handler) Process(msg *domain.Message, user *domain.User) error {
	u, err := h.getUser(user)
	if err != nil {
		return errors.Wrap(err, "Process() could not get user")
	}

	fmt.Println(u)

	return nil
}

func (h *Handler) getUser(user *domain.User) (*postgres.User, error) {
	// Get user from message
	u, err := h.us.Get(user.ID)
	if err != nil {
		return u, errors.Wrap(err, "Process() could not fetch user")
	}

	// User does not exist, create it.
	if u == nil {
		if err := h.us.Create(user.ID, user.Name); err != nil {
			return u, errors.Wrap(err, "Process() could not create user")
		}

		u, err := h.us.Get(user.ID)
		if err != nil || u == nil {
			return u, errors.Wrap(err, "Process() could not fetch user")
		}
	}

	return u, nil
}
