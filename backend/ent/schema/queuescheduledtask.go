package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// QueueScheduledTask holds the schema definition for the QueueScheduledTask entity.
type QueueScheduledTask struct {
	ent.Schema
}

// Fields of the QueueScheduledTask.
func (QueueScheduledTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.Int("interval"),
		field.Time("last_run_timestamp").Default(time.Now),
		field.String("status"),
	}
}

// Edges of the QueueScheduledTask.
func (QueueScheduledTask) Edges() []ent.Edge {
	return nil
}
