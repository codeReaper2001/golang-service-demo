package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Student holds the schema definition for the Student entity.
type Teacher struct {
	ent.Schema
}

// Fields of the Person.
func (Teacher) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("teacher_id"),
		field.String("teacher_name"),
	}
}

// Edges of the Person.
func (Teacher) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("students", Student.Type),
	}
}
