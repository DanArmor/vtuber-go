package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ReportedStream holds the schema definition for the ReportedStream entity.
type ReportedStream struct {
	ent.Schema
}

// Fields of the ReportedStream.
func (ReportedStream) Fields() []ent.Field {
	return []ent.Field{
		field.String("video_id"),
		field.Int("vtuber_id"),
		field.Time("available_at"),
	}
}

// Edges of the ReportedStream.
func (ReportedStream) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", Vtuber.Type).Ref("repoted_streams").Required().Unique(),
	}
}
