package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type Artifact struct{ ent.Schema }

func (Artifact) Fields() []ent.Field {
	return []ent.Field{
		field.Int("case_id").Optional(),
		field.Int("scan_id").Optional(),
		field.String("kind"),                // url|file|screenshot
		field.String("path_or_url"),        // abs path or URL
		field.String("sha256").Default(""),
		field.String("meta").Default("{}"),
		field.Time("added_at").Default(time.Now),
	}
}
