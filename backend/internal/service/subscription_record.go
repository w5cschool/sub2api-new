package service

import "time"

type SubscriptionRecord struct {
	ID             int64
	UserID         int64
	GroupID        int64
	SubscriptionID *int64
	PriceUSD       float64
	ValidityDays   int
	StartsAt       time.Time
	ExpiresAt      time.Time
	AssignedBy     *int64
	AssignedAt     time.Time
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	User           *User
	Group          *Group
	Subscription   *UserSubscription
	AssignedByUser *User
}

type SubscriptionRecordFilters struct {
	UserID    *int64
	GroupID   *int64
	StartTime *time.Time
	EndTime   *time.Time
}

type SubscriptionRecordStats struct {
	TotalAmountUSD float64
	RecordCount    int64
}
