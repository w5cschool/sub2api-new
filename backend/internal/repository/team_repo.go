package repository

import (
	"context"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dbteam "github.com/Wei-Shaw/sub2api/ent/team"
	dbteammember "github.com/Wei-Shaw/sub2api/ent/teammember"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	entsql "entgo.io/ent/dialect/sql"
)

type teamRepository struct {
	client *dbent.Client
}

func NewTeamRepository(client *dbent.Client) service.TeamRepository {
	return &teamRepository{client: client}
}

func (r *teamRepository) Create(ctx context.Context, team *service.Team) error {
	if team == nil {
		return service.ErrTeamNotFound
	}
	client := clientFromContext(ctx, r.client)
	builder := client.Team.Create().
		SetName(team.Name).
		SetStatus(team.Status)
	if team.Description != nil {
		builder.SetDescription(*team.Description)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, service.ErrTeamNameExists)
	}
	applyTeamEntityToService(team, created)
	return nil
}

func (r *teamRepository) Update(ctx context.Context, team *service.Team) error {
	if team == nil || team.ID <= 0 {
		return service.ErrTeamNotFound
	}
	client := clientFromContext(ctx, r.client)
	builder := client.Team.UpdateOneID(team.ID).
		SetName(team.Name).
		SetStatus(team.Status)
	if team.Description == nil {
		builder.ClearDescription()
	} else {
		builder.SetDescription(*team.Description)
	}
	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrTeamNotFound, service.ErrTeamNameExists)
	}
	applyTeamEntityToService(team, updated)
	return nil
}

func (r *teamRepository) GetByID(ctx context.Context, id int64) (*service.Team, error) {
	client := clientFromContext(ctx, r.client)
	entity, err := client.Team.Query().
		Where(dbteam.IDEQ(id)).
		WithMembers(func(q *dbent.TeamMemberQuery) {
			q.WithUser()
			q.Order(dbent.Desc(dbteammember.FieldRole), dbent.Asc(dbteammember.FieldUserID))
		}).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrTeamNotFound, nil)
	}
	return teamEntityToService(entity), nil
}

func (r *teamRepository) List(ctx context.Context, params pagination.PaginationParams, filters service.TeamFilters) ([]service.Team, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.Team.Query()
	q = applyTeamFilters(q, filters)

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	q = q.WithMembers(func(q *dbent.TeamMemberQuery) {
		q.WithUser()
		q.Order(dbent.Desc(dbteammember.FieldRole), dbent.Asc(dbteammember.FieldUserID))
	}).
		Offset(params.Offset()).
		Limit(params.Limit())
	for _, order := range teamListOrder(params) {
		q = q.Order(order)
	}
	entities, err := q.All(ctx)
	if err != nil {
		return nil, nil, err
	}
	return teamEntitiesToService(entities), paginationResultFromTotal(int64(total), params), nil
}

func (r *teamRepository) SoftDelete(ctx context.Context, id int64) error {
	if tx := dbent.TxFromContext(ctx); tx != nil {
		return r.softDelete(ctx, tx.Client(), id)
	}
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	txCtx := dbent.NewTxContext(ctx, tx)
	defer func() { _ = tx.Rollback() }()
	if err := r.softDelete(txCtx, tx.Client(), id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *teamRepository) ReplaceMembers(ctx context.Context, teamID int64, members []service.TeamMemberInput) error {
	if tx := dbent.TxFromContext(ctx); tx != nil {
		return r.replaceMembers(ctx, tx.Client(), teamID, members)
	}
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	txCtx := dbent.NewTxContext(ctx, tx)
	defer func() { _ = tx.Rollback() }()
	if err := r.replaceMembers(txCtx, tx.Client(), teamID, members); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *teamRepository) ListMembers(ctx context.Context, teamID int64) ([]service.TeamMember, error) {
	client := clientFromContext(ctx, r.client)
	entities, err := client.TeamMember.Query().
		Where(dbteammember.TeamIDEQ(teamID)).
		WithTeam().
		WithUser().
		Order(dbent.Desc(dbteammember.FieldRole), dbent.Asc(dbteammember.FieldUserID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return teamMemberEntitiesToService(entities), nil
}

func (r *teamRepository) GetMembershipByUserID(ctx context.Context, userID int64) (*service.TeamMember, error) {
	client := clientFromContext(ctx, r.client)
	entity, err := client.TeamMember.Query().
		Where(dbteammember.UserIDEQ(userID)).
		WithTeam().
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrTeamMembershipNotFound, nil)
	}
	return teamMemberEntityToService(entity), nil
}

func (r *teamRepository) softDelete(ctx context.Context, client *dbent.Client, id int64) error {
	if err := client.Team.DeleteOneID(id).Exec(ctx); err != nil {
		return translatePersistenceError(err, service.ErrTeamNotFound, nil)
	}
	_, err := client.TeamMember.Delete().Where(dbteammember.TeamIDEQ(id)).Exec(ctx)
	return err
}

func (r *teamRepository) replaceMembers(ctx context.Context, client *dbent.Client, teamID int64, members []service.TeamMemberInput) error {
	if _, err := client.Team.Query().Where(dbteam.IDEQ(teamID)).Only(ctx); err != nil {
		return translatePersistenceError(err, service.ErrTeamNotFound, nil)
	}
	if _, err := client.TeamMember.Delete().Where(dbteammember.TeamIDEQ(teamID)).Exec(ctx); err != nil {
		return err
	}
	if len(members) == 0 {
		return nil
	}
	builders := make([]*dbent.TeamMemberCreate, 0, len(members))
	for _, member := range members {
		builders = append(builders, client.TeamMember.Create().
			SetTeamID(teamID).
			SetUserID(member.UserID).
			SetRole(member.Role))
	}
	return client.TeamMember.CreateBulk(builders...).Exec(ctx)
}

func applyTeamFilters(q *dbent.TeamQuery, filters service.TeamFilters) *dbent.TeamQuery {
	if status := strings.TrimSpace(filters.Status); status != "" {
		q = q.Where(dbteam.StatusEQ(status))
	}
	if search := strings.TrimSpace(filters.Search); search != "" {
		q = q.Where(dbteam.Or(
			dbteam.NameContainsFold(search),
			dbteam.DescriptionContainsFold(search),
		))
	}
	return q
}

func teamListOrder(params pagination.PaginationParams) []func(*entsql.Selector) {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := params.NormalizedSortOrder(pagination.SortOrderDesc)
	field := dbteam.FieldID
	switch sortBy {
	case "name":
		field = dbteam.FieldName
	case "status":
		field = dbteam.FieldStatus
	case "created_at":
		field = dbteam.FieldCreatedAt
	case "updated_at":
		field = dbteam.FieldUpdatedAt
	}
	if sortOrder == pagination.SortOrderAsc {
		return []func(*entsql.Selector){dbent.Asc(field), dbent.Asc(dbteam.FieldID)}
	}
	return []func(*entsql.Selector){dbent.Desc(field), dbent.Desc(dbteam.FieldID)}
}

func teamEntityToService(entity *dbent.Team) *service.Team {
	if entity == nil {
		return nil
	}
	out := &service.Team{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}
	out.Members = teamMemberEntitiesToService(entity.Edges.Members)
	return out
}

func teamEntitiesToService(entities []*dbent.Team) []service.Team {
	out := make([]service.Team, 0, len(entities))
	for _, entity := range entities {
		if team := teamEntityToService(entity); team != nil {
			out = append(out, *team)
		}
	}
	return out
}

func teamMemberEntityToService(entity *dbent.TeamMember) *service.TeamMember {
	if entity == nil {
		return nil
	}
	out := &service.TeamMember{
		ID:        entity.ID,
		TeamID:    entity.TeamID,
		UserID:    entity.UserID,
		Role:      entity.Role,
		JoinedAt:  entity.JoinedAt,
		UpdatedAt: entity.UpdatedAt,
	}
	if entity.Edges.Team != nil {
		out.Team = teamEntityToService(entity.Edges.Team)
	}
	if entity.Edges.User != nil {
		out.User = userEntityToService(entity.Edges.User)
	}
	return out
}

func teamMemberEntitiesToService(entities []*dbent.TeamMember) []service.TeamMember {
	out := make([]service.TeamMember, 0, len(entities))
	for _, entity := range entities {
		if member := teamMemberEntityToService(entity); member != nil {
			out = append(out, *member)
		}
	}
	return out
}

func applyTeamEntityToService(dst *service.Team, src *dbent.Team) {
	if dst == nil || src == nil {
		return
	}
	dst.ID = src.ID
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.DeletedAt = src.DeletedAt
}
