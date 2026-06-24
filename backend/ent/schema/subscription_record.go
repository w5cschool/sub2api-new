package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SubscriptionRecord holds immutable admin subscription assignment records.
type SubscriptionRecord struct {
	ent.Schema
}

func (SubscriptionRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "subscription_records"},
	}
}

func (SubscriptionRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (SubscriptionRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.Int64("group_id"),
		field.Int64("subscription_id").
			Optional().
			Nillable(),
		field.Float("price_usd").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).
			Default(0),
		field.Int("validity_days").
			Default(30),
		field.Time("starts_at").
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("expires_at").
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Int64("assigned_by").
			Optional().
			Nillable(),
		field.Time("assigned_at").
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.String("notes").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}),
	}
}

func (SubscriptionRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("subscription_records").
			Field("user_id").
			Unique().
			Required(),
		edge.From("group", Group.Type).
			Ref("subscription_records").
			Field("group_id").
			Unique().
			Required(),
		edge.From("subscription", UserSubscription.Type).
			Ref("subscription_records").
			Field("subscription_id").
			Unique(),
		edge.From("assigned_by_user", User.Type).
			Ref("assigned_subscription_records").
			Field("assigned_by").
			Unique(),
	}
}

func (SubscriptionRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("group_id"),
		index.Fields("subscription_id"),
		index.Fields("assigned_by"),
		index.Fields("created_at"),
		index.Fields("user_id", "created_at"),
		index.Fields("group_id", "created_at"),
	}
}
