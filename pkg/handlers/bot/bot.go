package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
)

const (
	priceToBuyVideo = 1
)

type (
	Handler struct {
		debug        bool
		debugChannel string
		m            Messager
		rs           RankingStorage
		yt           YoutubeGetter
		balance      BalanceModifier
	}

	Messager interface {
		SendMessage(text string, channel string)
	}

	RankingStorage interface {
		GetRanking() ([]domain.User, error)
	}

	YoutubeGetter interface {
		GetVideo(q string) (string, error)
	}

	BalanceModifier interface {
		IncrementBalance(id string, inc int) error
	}
)

var (
	UnableToGetRankingError    = errors.New("could not get ranking")
	UnableToGetVideoError      = errors.New("could not get video")
	UnableToModifyBalanceError = errors.New("could not modify balance")

	acceptedBalance = map[string]bool{
		"tabela": true,
	}

	acceptedHelp = map[string]bool{
		"help":  true,
		"ajuda": true,
		"?":     true,
	}

	acceptedVideo = "youtube"
)

func NewHandler(debug bool, debugChannel string, m Messager, rs RankingStorage, yt YoutubeGetter, balance BalanceModifier) *Handler {
	return &Handler{
		debug:        debug,
		debugChannel: debugChannel,
		m:            m,
		rs:           rs,
		yt:           yt,
		balance:      balance,
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

	text := h.getFormattedMessage(prefix, msg.Content)
	if acceptedBalance[text] {
		return h.sendRanking(channel)
	}
	if acceptedHelp[text] {
		return h.sendHelpResponse(channel)
	}
	if strings.Contains(text, acceptedVideo) {
		return h.sendVideo(user, text, channel)
	}

	return h.sendNotUnderstood(user, channel)
}

func (h *Handler) getFormattedMessage(prefix, text string) string {
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	return text
}

func (h *Handler) sendNotUnderstood(user *domain.User, channel string) error {
	h.m.SendMessage(fmt.Sprintf("Desculpe, %s. Nao entendi o que voce quis dizer.\nCaso precise de ajude, basta enviar:\n```@Chicoin ajuda```", user.Name), channel)

	return nil
}

func (h *Handler) sendRanking(channel string) error {
	ranking, err := h.rs.GetRanking()
	if err != nil {
		return fmt.Errorf("%w : %v", UnableToGetRankingError, err)
	}

	var response string
	for i := 0; i < len(ranking); i++ {
		response += fmt.Sprintf("#%d ", i+1)
		response += "@" + ranking[i].Name
		response += ": "
		response += strconv.Itoa(ranking[i].Balance)
		response += domain.ChicoinEmoji
		response += "\n"
	}

	h.m.SendMessage(response, channel)

	return nil
}

func (h *Handler) sendHelpResponse(channel string) error {
	var response string

	response += "Comandos disponiveis:\n"
	response += "- tabela: Mostra o ranking atual\n\n"
	response += fmt.Sprintf("- youtube: Compra um video no youtube por %d %s\n\n", priceToBuyVideo, domain.ChicoinEmoji)
	response += "Novos comandos em breve!"

	h.m.SendMessage(response, channel)

	return nil
}

func (h *Handler) sendVideo(user *domain.User, query, channel string) error {
	query = strings.TrimPrefix(query, acceptedVideo)
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)

	if query == "" {
		h.m.SendMessage(fmt.Sprintf("Por apenas %d %s voce pode comprar um video!\nBasta escrever: \n```@Chicoin youtube palavra-chave```", priceToBuyVideo, domain.ChicoinEmoji), channel)

		return nil
	}

	if user.Balance < priceToBuyVideo {
		h.m.SendMessage(fmt.Sprintf("Desculpe, *%s*. Voce *nao* tem %s suficientes para comprar um video. :disappointed:", user.Name, domain.ChicoinEmoji), channel)

		return nil
	}

	err := h.balance.IncrementBalance(user.Id, -priceToBuyVideo)
	if err != nil {
		return fmt.Errorf("%w : %v", UnableToModifyBalanceError, err)
	}

	url, err := h.yt.GetVideo(query)
	if err != nil {
		return fmt.Errorf("%w : %v", UnableToGetVideoError, err)
	}

	h.m.SendMessage(fmt.Sprintf("*%s* comprou um video por *%d %s*.\n\n%s", user.Name, priceToBuyVideo, domain.ChicoinEmoji, url), channel)

	return nil
}
