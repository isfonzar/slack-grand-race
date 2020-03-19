package youtube

import (
	"fmt"
	"strings"

	yt "google.golang.org/api/youtube/v3"
)

type (
	Handler struct {
		youtubeService *yt.Service
	}
)

func NewHandler(youtubeService *yt.Service) *Handler {
	return &Handler{
		youtubeService: youtubeService,
	}
}

func (h *Handler) GetVideo(q string) (string, error) {
	call := h.youtubeService.Search.List("id,snippet").
		Q(q).Type("video").
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("error making search API call: %v", err)
	}

	if response.PageInfo.TotalResults <= 0 {
		return fmt.Sprintf("Nenhum video encontrado para a busca: `%s`", q), nil
	}

	var msgList []string
	for _, item := range response.Items {
		if item.Id.VideoId != "" {
			msgList = append(msgList, fmt.Sprintf("%s - http://youtu.be/%s", item.Snippet.Title, item.Id.VideoId))
		}
	}

	return strings.Join(msgList, "\n"), nil
}
