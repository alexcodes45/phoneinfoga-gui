package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Case represents a grouping of scans and artifacts for a client or investigation.
type Case struct {
	ent.Schema
}

func (Case) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("client_ref").Optional().Nillable(),
		field.Text("notes").Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Case) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("scans", Scan.Type),
		edge.To("artifacts", Artifact.Type),
	}
}
