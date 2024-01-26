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
	MembershipLevelId   string
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
	MembershipLevelId   string
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
