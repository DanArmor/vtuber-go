package types

type Vtuber struct {
	YoutubeChannelId string   `json:"id,omitempty"`
	ChannelName      string   `json:"name,omitempty"`
	EnglishName      string   `json:"english_name,omitempty"`
	PhotoURL         string   `json:"photo,omitempty"`
	Twitter          string   `json:"twitter,omitempty"`
	VideoCount       int      `json:"video_count,omitempty"`
	SubscriberCount  int      `json:"subscriber_count,omitempty"`
	ClipCount        int      `json:"clip_count,omitempty"`
	TopTopics        []string `json:"top_topics,omitempty"`
	Inactive         bool     `json:"inactive,omitempty"`
	Twitch           string   `json:"twitch,omitempty"`
	BannerURL        string   `json:"banner,omitempty"`
	Description      string   `json:"description,omitempty"`
	WaveName         string   `json:"group,omitempty"`
	CompanyName      string   `json:"org,omitempty"`
}
