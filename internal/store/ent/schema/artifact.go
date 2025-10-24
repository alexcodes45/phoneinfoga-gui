package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Artifact tracks evidence attached to scans or cases.
type Artifact struct {
	ent.Schema
}

func (Artifact) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").NotEmpty(),
		field.String("path_or_url").NotEmpty(),
		field.String("sha256").NotEmpty(),
		field.JSON("meta", map[string]any{}).Optional().Nillable(),
		field.Time("added_at").Default(time.Now),
	}
}

func (Artifact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("scan", Scan.Type).
			Ref("artifacts").
			Unique().
			Optional(),
		edge.From("case", Case.Type).
			Ref("artifacts").
			Unique().
			Optional(),
	}
}

func (Artifact) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sha256").Unique(),
	}
}
