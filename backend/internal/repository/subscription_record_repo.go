package repository

import (
	"context"
	"fmt"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/subscriptionrecord"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type subscriptionRecordRepository struct {
	client *dbent.Client
}

func NewSubscriptionRecordRepository(client *dbent.Client) service.SubscriptionRecordRepository {
	return &subscriptionRecordRepository{client: client}
}

func (r *subscriptionRecordRepository) Create(ctx context.Context, record *service.SubscriptionRecord) error {
	if record == nil {
		return service.ErrSubscriptionNilInput
	}

	client := clientFromContext(ctx, r.client)
	builder := client.SubscriptionRecord.Create().
		SetUserID(record.UserID).
		SetGroupID(record.GroupID).
		SetNillableSubscriptionID(record.SubscriptionID).
		SetOperation(normalizeSubscriptionRecordOperation(record.Operation)).
		SetPriceUsd(record.PriceUSD).
		SetValidityDays(record.ValidityDays).
		SetStartsAt(record.StartsAt).
		SetExpiresAt(record.ExpiresAt).
		SetNillableAssignedBy(record.AssignedBy).
		SetAssignedAt(record.AssignedAt).
		SetNotes(record.Notes)

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	applySubscriptionRecordEntityToService(record, created)
	return nil
}

func (r *subscriptionRecordRepository) List(ctx context.Context, params pagination.PaginationParams, filters service.SubscriptionRecordFilters) ([]service.SubscriptionRecord, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := r.applyFilters(client.SubscriptionRecord.Query(), filters)

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	records, err := q.
		WithUser().
		WithGroup().
		WithSubscription(func(q *dbent.UserSubscriptionQuery) {
			q.WithUser().WithGroup().WithAssignedByUser()
		}).
		WithAssignedByUser().
		Order(dbent.Desc(subscriptionrecord.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	return subscriptionRecordEntitiesToService(records), paginationResultFromTotal(int64(total), params), nil
}

func (r *subscriptionRecordRepository) ListForExport(ctx context.Context, filters service.SubscriptionRecordFilters, limit int) ([]service.SubscriptionRecord, error) {
	client := clientFromContext(ctx, r.client)
	q := r.applyFilters(client.SubscriptionRecord.Query(), filters).
		WithUser().
		WithGroup().
		WithSubscription(func(q *dbent.UserSubscriptionQuery) {
			q.WithUser().WithGroup().WithAssignedByUser()
		}).
		WithAssignedByUser().
		Order(dbent.Desc(subscriptionrecord.FieldCreatedAt))
	if limit > 0 {
		q = q.Limit(limit)
	}

	records, err := q.All(ctx)
	if err != nil {
		return nil, err
	}
	return subscriptionRecordEntitiesToService(records), nil
}

func (r *subscriptionRecordRepository) Stats(ctx context.Context, filters service.SubscriptionRecordFilters) (*service.SubscriptionRecordStats, error) {
	client := clientFromContext(ctx, r.client)

	clauses := []string{"1 = 1"}
	args := make([]any, 0, 4)
	addArg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.UserID != nil {
		clauses = append(clauses, "user_id = "+addArg(*filters.UserID))
	}
	if filters.GroupID != nil {
		clauses = append(clauses, "group_id = "+addArg(*filters.GroupID))
	}
	if filters.StartTime != nil {
		clauses = append(clauses, "created_at >= "+addArg(*filters.StartTime))
	}
	if filters.EndTime != nil {
		clauses = append(clauses, "created_at <= "+addArg(*filters.EndTime))
	}

	query := `SELECT COALESCE(SUM(price_usd), 0), COUNT(*) FROM subscription_records WHERE ` + strings.Join(clauses, " AND ")
	rows, err := client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	stats := &service.SubscriptionRecordStats{}
	if rows.Next() {
		if err := rows.Scan(&stats.TotalAmountUSD, &stats.RecordCount); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *subscriptionRecordRepository) applyFilters(q *dbent.SubscriptionRecordQuery, filters service.SubscriptionRecordFilters) *dbent.SubscriptionRecordQuery {
	if filters.UserID != nil {
		q = q.Where(subscriptionrecord.UserIDEQ(*filters.UserID))
	}
	if filters.GroupID != nil {
		q = q.Where(subscriptionrecord.GroupIDEQ(*filters.GroupID))
	}
	if filters.StartTime != nil {
		q = q.Where(subscriptionrecord.CreatedAtGTE(*filters.StartTime))
	}
	if filters.EndTime != nil {
		q = q.Where(subscriptionrecord.CreatedAtLTE(*filters.EndTime))
	}
	return q
}

func subscriptionRecordEntityToService(m *dbent.SubscriptionRecord) *service.SubscriptionRecord {
	if m == nil {
		return nil
	}
	out := &service.SubscriptionRecord{
		ID:             m.ID,
		UserID:         m.UserID,
		GroupID:        m.GroupID,
		SubscriptionID: m.SubscriptionID,
		Operation:      normalizeSubscriptionRecordOperation(m.Operation),
		PriceUSD:       m.PriceUsd,
		ValidityDays:   m.ValidityDays,
		StartsAt:       m.StartsAt,
		ExpiresAt:      m.ExpiresAt,
		AssignedBy:     m.AssignedBy,
		AssignedAt:     m.AssignedAt,
		Notes:          derefString(m.Notes),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
	}
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
	}
	if m.Edges.Subscription != nil {
		out.Subscription = userSubscriptionEntityToService(m.Edges.Subscription)
	}
	if m.Edges.AssignedByUser != nil {
		out.AssignedByUser = userEntityToService(m.Edges.AssignedByUser)
	}
	return out
}

func normalizeSubscriptionRecordOperation(operation string) string {
	if operation == "" {
		return service.SubscriptionRecordOperationAssign
	}
	return operation
}

func subscriptionRecordEntitiesToService(models []*dbent.SubscriptionRecord) []service.SubscriptionRecord {
	out := make([]service.SubscriptionRecord, 0, len(models))
	for i := range models {
		if r := subscriptionRecordEntityToService(models[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

func applySubscriptionRecordEntityToService(dst *service.SubscriptionRecord, src *dbent.SubscriptionRecord) {
	if dst == nil || src == nil {
		return
	}
	dst.ID = src.ID
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}
