package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vtuber holds the schema definition for the Vtuber entity.
type Vtuber struct {
	ent.Schema
}

// Fields of the Vtuber.
func (Vtuber) Fields() []ent.Field {
	return []ent.Field{
		field.String("channel_name"),
		field.String("youtube_channel_id"),
		field.String("english_name").Optional(),
		field.String("photo_url").Optional(),
		field.String("banner_url").Optional(),
		field.String("twitter").Optional(),
		field.Int("video_count").Optional(),
		field.Int("subscriber_count").Optional(),
		field.Int("clip_count").Optional(),
		field.Strings("top_topics").Optional(),
		field.String("description").Optional(),
		field.Bool("inactive"),
		field.String("twitch").Optional(),
	}
}

// Edges of the Vtuber.
func (Vtuber) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wave", Wave.Type).Ref("vtubers").Required().Unique(),
		edge.From("users", User.Type).Ref("vtubers"),
	}
}
