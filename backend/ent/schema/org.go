package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Org holds the schema definition for the Org entity.
type Org struct {
	ent.Schema
}

// Fields of the Org.
func (Org) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the Org.
func (Org) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("waves", Wave.Type),
	}
}
