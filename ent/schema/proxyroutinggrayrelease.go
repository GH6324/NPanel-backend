package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingGrayRelease struct{ ent.Schema }

func (ProxyRoutingGrayRelease) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_gray_release"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingGrayRelease) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Routing gray release ID"),
		field.String("profile_code").MaxLen(128).NotEmpty().Comment("Routing profile code"),
		field.String("name").MaxLen(255).NotEmpty().Comment("Gray release name"),
		field.String("status").MaxLen(32).Default("draft").Comment("draft/running/paused/completed/rolled_back"),
		field.Int("batch_no").Default(0).Comment("Current gray batch number"),
		field.String("target_type").MaxLen(32).Default("user").Comment("user/user_subscribe/subscribe/node"),
		field.Text("target_ids_json").Default("[]").Comment("Target IDs JSON array"),
		field.String("operator").MaxLen(128).Default("").Comment("Operator"),
		field.Text("rollback_reason").Optional().Comment("Rollback reason"),
		field.Time("started_at").Optional().Comment("Started at"),
		field.Time("ended_at").Optional().Comment("Ended at"),
		field.Text("release_json").Default("{}").Comment("Release metadata JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingGrayRelease) Edges() []ent.Edge { return nil }
