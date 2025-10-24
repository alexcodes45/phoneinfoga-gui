package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type Case struct{ ent.Schema }

func (Case) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("client_ref").Default(""),
		field.String("notes").Default(""),
		field.Time("created_at").Default(time.Now),
	}
}
