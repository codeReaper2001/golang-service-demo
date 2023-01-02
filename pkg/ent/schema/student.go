package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Student holds the schema definition for the Student entity.
type Student struct {
	ent.Schema
}

// Fields of the Person.
func (Student) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("stu_id"),
		field.String("name"),
		field.Int32("age"),
	}
}

// Edges of the Person.
func (Student) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("teacher", Teacher.Type).
			Ref("students").
			Unique(),
	}
}
