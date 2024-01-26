package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"log/slog"
)

// CreateMembershipHandler is an abstraction for MemberShipStarted use-case handler.
type CreateMembershipHandler interface {
	HandleRequest(ctx context.Context, request CreateMemberShipRequest) error
}

type NewCreateMembershipArgs struct {
	Logger     *slog.Logger
	Config     *config.Config
	Repository domain.Repository
}

func NewCreateMembership(args NewCreateMembershipArgs) (*CreateMembership, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	return &CreateMembership{
		logger:     args.Logger,
		config:     args.Config,
		repository: args.Repository,
	}, nil
}

// CreateMembership is a request handler with all dependencies initialized.
type CreateMembership struct {
	logger     *slog.Logger
	config     *config.Config
	repository domain.Repository
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

func (h *CreateMembership) HandleRequest(ctx context.Context, request CreateMemberShipRequest) error {
	logger := h.logger
	repository := h.repository
	errorResponse := errors.New("")

	logger.Info("request at handler level", "request", request)

	err := repository.CreateMembership(ctx, domain.CreateMembershipArgs{
		Type:                request.Type,
		LiveMode:            request.LiveMode,
		Attempt:             request.Attempt,
		Created:             request.Created,
		EventId:             request.EventId,
		Id:                  request.Id,
		Amount:              request.Amount,
		Object:              request.Object,
		Paused:              request.Paused,
		Status:              request.Status,
		Canceled:            request.Canceled,
		Currency:            request.Currency,
		PspId:               request.PspId,
		MembershipLevelId:   request.MembershipLevelId,
		MembershipLevelName: request.MembershipLevelName,
		StartedAt:           request.StartedAt,
		CanceledAt:          request.CanceledAt,
		NoteHidden:          request.NoteHidden,
		SupportNote:         request.SupportNote,
		SupporterName:       request.SupporterName,
		SupporterId:         request.SupporterId,
		SupporterEmail:      request.SupporterEmail,
		CurrentPeriodEnd:    request.CurrentPeriodEnd,
		CurrentPeriodStart:  request.CurrentPeriodStart,
	})
	if err != nil {
		logger.Error("error on repository.create_membership", "request", request, "err", err)
		return errorResponse
	}

	return nil
}
