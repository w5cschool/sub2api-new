package service

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type SubscriptionRecordRepository interface {
	Create(ctx context.Context, record *SubscriptionRecord) error
	List(ctx context.Context, params pagination.PaginationParams, filters SubscriptionRecordFilters) ([]SubscriptionRecord, *pagination.PaginationResult, error)
	Stats(ctx context.Context, filters SubscriptionRecordFilters) (*SubscriptionRecordStats, error)
}
