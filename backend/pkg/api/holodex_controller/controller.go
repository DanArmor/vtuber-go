package holodexcontroller

import (
	"context"
	"net/http"
	"strings"

	"github.com/DanArmor/go-holodex"
	"github.com/DanArmor/vtuber-go/ent"
	"go.uber.org/zap"
)

type HolodexController struct {
	apiKey string
	client *holodex.APIClient
	logger *zap.Logger
}

func NewHolodexController(apiKey string, logger *zap.Logger) *HolodexController {
	newLogger := logger.With(zap.String("Service", "HolodexController"))
	newLogger.Debug("Constructing HolodexController")

	holodexConfig := holodex.NewConfiguration()
	holodexConfig.DefaultHeader["X-APIKEY"] = apiKey
	holodexConfig.UserAgent = "Vtuber-Go"

	return &HolodexController{
		apiKey: apiKey,
		client: holodex.NewAPIClient(holodexConfig),
		logger: newLogger,
	}
}

func (hc *HolodexController) GetVtubersUpcomingVideos(ctx context.Context, vtubers []*ent.Vtuber) []holodex.Video {
	hc.logger.Debug("Getting vtubers upcoming videos")
	ids := make([]string, 0, len(vtubers))
	for i := range vtubers {
		if vtubers[i].YoutubeChannelID != "" {
			hc.logger.Debug(
				"Added vtuber to the list",
				zap.String("ChannelName", vtubers[i].ChannelName),
				zap.String("ChannelId", vtubers[i].YoutubeChannelID),
			)
			ids = append(ids, vtubers[i].YoutubeChannelID)
		}
	}
	queryParam := strings.Join(ids, ",")
	hc.logger.Debug("queryParam", zap.String("param", queryParam))

	request := hc.client.DefaultApi.GetCachedLive(ctx).Channels(queryParam)

	videos, response, err := hc.client.DefaultApi.GetCachedLiveExecute(request)
	for _, video := range videos {
		if video.Channel.Id == nil {
			hc.logger.Warn("Nil channel id", zap.String("VideoId", *video.Id))
			continue
		}
		hc.logger.Debug("Video info", zap.String("VideoId", *video.Id), zap.String("ChannelId", *video.Channel.Id))
	}
	defer response.Body.Close()
	if err != nil || response.StatusCode != http.StatusOK {
		hc.logger.Error("Notify error", zap.Error(err), zap.Int("StatusCode", response.StatusCode))
		return []holodex.Video{}
	}
	return videos
}
