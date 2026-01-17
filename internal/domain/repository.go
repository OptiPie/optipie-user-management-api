package domain

import (
	"context"
	"github.com/OptiPie/optipie-user-management-api/internal/domain/models"
	"time"
)

type Repository interface {
	CreateMembership(ctx context.Context, args CreateMembershipArgs) error
	GetMembershipByEmail(ctx context.Context, email string) (models.Membership, error)
	UpdateMembershipByEmail(ctx context.Context, email string, args UpdateMembershipArgs) error
	DeleteMembershipByEmail(ctx context.Context, email string) error
	CreateAnalytics(ctx context.Context, args CreateAnalyticsArgs) error
}

// CreateMembershipArgs to call CreateMembership repository method
type CreateMembershipArgs struct {
	Type                string
	LiveMode            bool
	Attempt             int32
	Created             time.Time
	EventId             int64
	Id                  int64
	Amount              float64
	Object              string
	Paused              string
	Status              string
	Canceled            string
	Currency            string
	PspId               string
	MembershipLevelId   int64
	MembershipLevelName string
	StartedAt           time.Time
	CanceledAt          time.Time
	NoteHidden          bool
	SupportNote         string
	SupporterName       string
	SupporterId         int64
	SupporterEmail      string
	CurrentPeriodEnd    time.Time
	CurrentPeriodStart  time.Time
}

type UpdateMembershipArgs struct {
	updated             time.Time
	Paused              string
	Status              string
	Canceled            string
	Currency            string
	PspId               string
	MembershipLevelId   int64
	MembershipLevelName string
	StartedAt           time.Time
	CanceledAt          time.Time
	NoteHidden          bool
	SupportNote         string
	SupporterName       string
	CurrentPeriodEnd    time.Time
	SupporterFeedback   string
	CancelAtPeriodEnd   string
	CurrentPeriodStart  time.Time
}

// DeleteMembershipArgs in case of soft deletion needed, this will be used
type DeleteMembershipArgs struct {
	updated            time.Time
	Status             string
	Canceled           string
	CanceledAt         time.Time
	CurrentPeriodEnd   time.Time
	SupporterFeedback  string
	CurrentPeriodStart time.Time
}

// CreateAnalyticsArgs to call CreateAnalytics repository method
type CreateAnalyticsArgs struct {
	Timestamp         int64
	StrategyName      string
	StrategySymbol    string
	StrategyPeriod    string
	StrategyDateRange string
}
