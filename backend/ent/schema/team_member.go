package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TeamMember holds the schema definition for team memberships.
type TeamMember struct {
	ent.Schema
}

func (TeamMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "team_members"},
	}
}

func (TeamMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("team_id"),
		field.Int64("user_id"),
		field.String("role").
			MaxLen(20).
			Default("member"),
		field.Time("joined_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (TeamMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("team", Team.Type).
			Ref("members").
			Field("team_id").
			Required().
			Unique(),
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id"),
	}
}

func (TeamMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("team_id"),
		index.Fields("user_id").Unique(),
		index.Fields("team_id", "role"),
	}
}
