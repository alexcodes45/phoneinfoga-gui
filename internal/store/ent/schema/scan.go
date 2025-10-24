package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type Scan struct{ ent.Schema }

func (Scan) Fields() []ent.Field {
	return []ent.Field{
		field.String("number_e164"),
		field.String("backend").Default("serve"),
		field.String("status").Default("queued"),
		field.Time("started_at").Default(time.Now),
		field.Time("finished_at").Nillable().Optional(),
		field.Int("duration_ms").Default(0),
		field.String("request_opts").Default("{}"),
		field.String("phoneinfoga_version").Default(""),
		field.Bytes("raw_json").Optional(),
		field.String("result_hash").Default(""),
	}
}
