package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Result stores structured data parsed from a PhoneInfoga response.
type Result struct {
	ent.Schema
}

func (Result) Fields() []ent.Field {
	return []ent.Field{
		field.String("country").Optional().Nillable(),
		field.String("region").Optional().Nillable(),
		field.String("carrier").Optional().Nillable(),
		field.String("line_type").Optional().Nillable(),
		field.JSON("reputation", map[string]any{}).Optional().Nillable(),
		field.JSON("sources", []map[string]any{}).Optional().Nillable(),
		field.Text("notes").Optional().Nillable(),
	}
}

func (Result) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("scan", Scan.Type).
			Ref("results").
			Unique().
			Required(),
	}
}
