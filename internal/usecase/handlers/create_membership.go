package handlers

import (
	"context"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"log/slog"
)

// CreateMembershipHandler is an abstraction for MemberShipStarted use-case handler.
type CreateMembershipHandler interface {
	HandleRequest(ctx context.Context, request CreateMemberShipRequest) error
}

type NewCreateMembershipArgs struct {
	Logger *slog.Logger
	Config *config.Config
}

func NewCreateMembership(args NewCreateMembershipArgs) (*CreateMembership, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	return &CreateMembership{
		logger: args.Logger,
		config: args.Config,
	}, nil
}

// CreateMembership is a request handler with all dependencies initialized.
type CreateMembership struct {
	logger *slog.Logger
	config *config.Config
}

// CreateMemberShipRequest represents necessary POST /api/v1/user/membership request data for handler.
type CreateMemberShipRequest struct {
	Type                string
	LiveMode            bool
	Attempt             int32
	Created             int64
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
	StartedAt           int64
	CanceledAt          int64
	NoteHidden          bool
	SupportNote         string
	SupporterName       string
	SupporterId         int64
	SupporterEmail      string
	CurrentPeriodEnd    int64
	CurrentPeriodStart  int64
}

type CreateMemberShipResponse struct {
}

func (h *CreateMembership) HandleRequest(ctx context.Context, request CreateMemberShipRequest) error {
	logger := h.logger
	logger.Info("Here is the request at handler level", request)

	// add dynamodb logic here.

	return nil
}
