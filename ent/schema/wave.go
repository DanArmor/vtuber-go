package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Wave holds the schema definition for the Wave entity.
type Wave struct {
	ent.Schema
}

// Fields of the Wave.
func (Wave) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the Wave.
func (Wave) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vtubers", Vtuber.Type),
		edge.From("org", Org.Type).Ref("waves").Required().Unique(),
	}
}
