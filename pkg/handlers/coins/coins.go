package coins

import (
	"errors"
	"fmt"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
)

type (
	Handler struct {
		debug bool
		us    Storage
		re    Reactioner
	}

	Reactioner interface {
		AddReaction(msg *domain.Message, reaction domain.Reaction) error
	}

	Storage interface {
		IncrementBalance(id string, inc int) error
	}
)

var (
	CouldNotChangeBalance    = errors.New("could not change balance")
	CouldNotAddReactionError = errors.New("could not add reaction")
)

func NewHandler(debug bool, us Storage, re Reactioner) *Handler {
	return &Handler{
		debug: debug,
		us:    us,
		re:    re,
	}
}

func (h *Handler) Give(msg *domain.Message, amount int) error {
	if err := h.us.IncrementBalance(msg.User, amount); err != nil {
		return fmt.Errorf("%w : %v", CouldNotChangeBalance, err)
	}

	// If debug is set to true, do not send reactions
	if h.debug {
		return nil
	}

	if err := h.re.AddReaction(msg, domain.ChicoinReaction); err != nil {
		return fmt.Errorf("%w : %v", CouldNotAddReactionError, err)
	}

	return nil
}
