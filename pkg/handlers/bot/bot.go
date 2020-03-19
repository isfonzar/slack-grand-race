package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
)

type (
	Handler struct {
		debug        bool
		debugChannel string
		m            Messager
		rs           RankingStorage
	}

	Messager interface {
		SendMessage(text string, channel string)
	}

	RankingStorage interface {
		GetRanking() ([]domain.User, error)
	}
)

var (
	UnableToGetRankingError = errors.New("could not get ranking")

	acceptedBalance = map[string]bool{
		"tabela": true,
	}

	acceptedHelp = map[string]bool{
		"help":  true,
		"ajuda": true,
		"?":     true,
	}
)

func NewHandler(debug bool, debugChannel string, m Messager, rs RankingStorage) *Handler {
	return &Handler{
		debug:        debug,
		debugChannel: debugChannel,
		m:            m,
		rs:           rs,
	}
}

func (h *Handler) Process(selfId string, msg *domain.Message, user *domain.User) error {
	prefix := fmt.Sprintf("<@%s> ", selfId)

	// Message coming from self
	if user.Id == selfId {
		return nil
	}

	// Is not calling the bot
	if !strings.HasPrefix(msg.Content, prefix) {
		return nil
	}

	channel := msg.Channel
	if h.debug && h.debugChannel != "" {
		channel = h.debugChannel
	}

	text := msg.Content
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	if acceptedBalance[text] {
		ranking, err := h.rs.GetRanking()
		if err != nil {
			return fmt.Errorf("%w : %v", UnableToGetRankingError, err)
		}

		var response string
		for i := 0; i < len(ranking); i++ {
			response += "@" + ranking[i].Name
			response += ": "
			response += strconv.Itoa(ranking[i].Balance)
			response += domain.ChicoinEmoji
			response += "\n"
		}

		h.m.SendMessage(response, channel)

		return nil
	}
	if acceptedHelp[text] {
		var response string

		response += "Comandos disponiveis:\n"
		response += "- tabela: Mostra o ranking atual\n\n"
		response += "Novos comandos em breve!"

		h.m.SendMessage(response, channel)

		return nil
	}

	return nil
}
