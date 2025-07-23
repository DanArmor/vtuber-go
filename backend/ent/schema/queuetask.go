package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// QueueTask holds the schema definition for the QueueTask entity.
type QueueTask struct {
	ent.Schema
}

// Fields of the QueueTask.
func (QueueTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("task_name"),
		field.JSON("data", map[string]any{}),
		field.String("status"),
		field.String("worker_name").Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the QueueTask.
func (QueueTask) Edges() []ent.Edge {
	return nil
}
