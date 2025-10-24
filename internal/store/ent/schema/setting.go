package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Setting stores non-sensitive configuration key/value pairs.
type Setting struct {
	ent.Schema
}

func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").NotEmpty().Immutable(),
		field.JSON("value", map[string]any{}).
			Optional().
			Nillable(),
	}
}

func (Setting) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key").Unique(),
	}
}
