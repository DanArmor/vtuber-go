package types

type Vtuber struct {
	ChannelName     string   `json:"channel_name,omitempty"`
	EnglishName     string   `json:"english_name,omitempty"`
	PhotoURL        string   `json:"photo_url,omitempty"`
	Twitter         string   `json:"twitter,omitempty"`
	VideoCount      int      `json:"video_count,omitempty"`
	SubscriberCount int      `json:"subscriber_count,omitempty"`
	ClipCount       int      `json:"clip_count,omitempty"`
	TopTopics       []string `json:"top_topics,omitempty"`
	Inactive        bool     `json:"inactive,omitempty"`
	Twitch          string   `json:"twitch,omitempty"`
	WaveName        string   `json:"wave_name,omitempty"`
	CompanyName     string   `json:"company_name,omitempty"`
}
