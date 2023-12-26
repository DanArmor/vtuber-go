package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tg_id").Unique(),
		field.String("first_name"),
		field.String("last_name").Optional(),
		field.String("username").Optional(),
		field.String("language_code").Optional(),
		field.String("photo_url").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vtubers", Vtuber.Type),
	}
}
