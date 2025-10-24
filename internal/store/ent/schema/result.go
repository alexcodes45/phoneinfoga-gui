package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Result struct{ ent.Schema }

func (Result) Fields() []ent.Field {
	return []ent.Field{
		field.Int("scan_id"),
		field.String("country").Default(""),
		field.String("region").Default(""),
		field.String("carrier").Default(""),
		field.String("line_type").Default(""),
		field.String("reputation").Default("{}"),
		field.String("sources").Default("[]"),
		field.String("notes").Default(""),
	}
}
