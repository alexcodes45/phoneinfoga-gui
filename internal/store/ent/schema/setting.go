package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Setting struct{ ent.Schema }

func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.String("key"),
		field.String("value"),
	}
}
