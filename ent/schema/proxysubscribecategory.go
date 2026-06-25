package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ProxySubscribeCategory holds the schema definition for the ProxySubscribeCategory entity.
type ProxySubscribeCategory struct {
	ent.Schema
}

func (ProxySubscribeCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "subscribe_category"},
		entsql.WithComments(true),
	}
}

func (ProxySubscribeCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive().Comment("商品分类ID"),
		field.Int64("parent_id").Default(0).Comment("父分类ID"),
		field.String("name").MaxLen(255).Default("").Comment("分类名称"),
		field.Text("description").Optional().Nillable().Comment("分类描述"),
		field.String("language").MaxLen(255).Default("").Comment("语言"),
		field.Bool("show").Default(true).Comment("是否显示"),
		field.Int32("sort").Default(0).Comment("排序"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
	}
}

func (ProxySubscribeCategory) Edges() []ent.Edge {
	return nil
}
