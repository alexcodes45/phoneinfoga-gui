package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Scan holds metadata about an executed PhoneInfoga scan.
type Scan struct {
	ent.Schema
}

func (Scan) Fields() []ent.Field {
	return []ent.Field{
		field.String("number_e164").NotEmpty(),
		field.String("backend").Default("serve"),
		field.String("status").Default("queued"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("started_at").Optional().Nillable(),
		field.Time("finished_at").Optional().Nillable(),
		field.Int("duration_ms").Default(0),
		field.JSON("request_opts", map[string]any{}).Optional().Nillable(),
		field.String("phoneinfoga_version").Optional().Nillable(),
		field.Bytes("raw_json").Optional().Nillable(),
		field.String("result_hash").Optional().Nillable(),
	}
}

func (Scan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("case", Case.Type).
			Ref("scans").
			Unique().
			Optional(),
		edge.To("results", Result.Type),
		edge.To("artifacts", Artifact.Type),
	}
}

func (Scan) Annotations() []schema.Annotation {
	return nil
}

func (Scan) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("number_e164", "created_at"),
	}
}
