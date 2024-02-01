package handlers

import (
	"context"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"log/slog"
)

// UpdateMembershipHandler is an abstraction for MembershipUpdated use-case handler.
type UpdateMembershipHandler interface {
	HandleRequest(ctx context.Context, request UpdateMembershipRequest) error
}

type NewUpdateMembershipArgs struct {
	Logger     *slog.Logger
	Config     *config.Config
	Repository domain.Repository
}

func NewUpdateMembership(args NewUpdateMembershipArgs) (*UpdateMembership, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if args.Repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	return &UpdateMembership{
		logger:     args.Logger,
		config:     args.Config,
		repository: args.Repository,
	}, nil
}

// UpdateMembership is a request handler with all dependencies initialized.
type UpdateMembership struct {
	logger     *slog.Logger
	config     *config.Config
	repository domain.Repository
}

// UpdateMembershipRequest represents necessary POST /api/v1/user/membership/update request data for handler.
type UpdateMembershipRequest struct {
	Paused              string
	Status              string
	Canceled            string
	Currency            string
	PspId               string
	MembershipLevelId   int64
	MembershipLevelName string
	StartedAt           int64
	CanceledAt          int64
	NoteHidden          bool
	SupportNote         string
	SupporterName       string
	SupporterEmail      string
	CurrentPeriodEnd    int64
	SupporterFeedback   string
	CancelAtPeriodEnd   string
	CurrentPeriodStart  int64
}

func (h *UpdateMembership) HandleRequest(ctx context.Context, request UpdateMembershipRequest) error {
	logger := h.logger
	repository := h.repository

	err := repository.UpdateMembershipByEmail(ctx, request.SupporterEmail, domain.UpdateMembershipArgs{
		Paused:              request.Paused,
		Status:              request.Status,
		Canceled:            request.Canceled,
		Currency:            request.Currency,
		PspId:               request.PspId,
		MembershipLevelId:   request.MembershipLevelId,
		MembershipLevelName: request.MembershipLevelName,
		StartedAt:           convertUnixToUTCTime(request.StartedAt),
		CanceledAt:          convertUnixToUTCTime(request.CanceledAt),
		NoteHidden:          request.NoteHidden,
		SupportNote:         request.SupportNote,
		SupporterName:       request.SupporterName,
		CurrentPeriodEnd:    convertUnixToUTCTime(request.CurrentPeriodEnd),
		SupporterFeedback:   request.SupporterFeedback,
		CancelAtPeriodEnd:   request.CancelAtPeriodEnd,
		CurrentPeriodStart:  convertUnixToUTCTime(request.CurrentPeriodStart),
	})
	if err != nil {
		logger.Error("error on repository.update_membership_by_email", "request", request, "err", err)
		return err
	}

	return nil
}
