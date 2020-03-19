package message

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
)

type (
	Handler struct {
		cg  CoinGiver
		log Logger
	}

	CoinGiver interface {
		Give(msg *domain.Message, amount int) error
	}

	Logger interface {
		Infow(msg string, keysAndValues ...interface{})
	}
)

var (
	CouldNotGiveCoinError = errors.New("could not give coin")
)

func NewHandler(cg CoinGiver, l Logger) *Handler {
	return &Handler{
		cg:  cg,
		log: l,
	}
}

func (h *Handler) Process(msg *domain.Message, user *domain.User) error {
	if msg == nil ||
		user == nil {
		return nil
	}

	h.log.Infow("processing message",
		"message", msg,
		"user", user,
	)

	rand.Seed(time.Now().UnixNano())

	v := rand.Intn(100)
	if v == 0 {
		h.log.Infow("giving coin by chance to user",
			"user", user,
		)
		if err := h.cg.Give(msg, 1); err != nil {
			return fmt.Errorf("%w : %v", CouldNotGiveCoinError, err)
		}
	}

	return nil
}
