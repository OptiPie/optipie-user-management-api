package handlers

import (
	"context"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"log/slog"
	"time"
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
	if args.Repository == nil {
		return nil, fmt.Errorf("repository is required")
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

// CreateMemberShipRequest represents necessary POST /api/v1/user/membership/create  request data for handler.
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
	MembershipLevelId   int64
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

	err := repository.CreateMembership(ctx, domain.CreateMembershipArgs{
		Type:                request.Type,
		LiveMode:            request.LiveMode,
		Attempt:             request.Attempt,
		Created:             convertUnixToUTCTime(request.Created),
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
		StartedAt:           convertUnixToUTCTime(request.StartedAt),
		CanceledAt:          convertUnixToUTCTime(request.CanceledAt),
		NoteHidden:          request.NoteHidden,
		SupportNote:         request.SupportNote,
		SupporterName:       request.SupporterName,
		SupporterId:         request.SupporterId,
		SupporterEmail:      request.SupporterEmail,
		CurrentPeriodEnd:    convertUnixToUTCTime(request.CurrentPeriodEnd),
		CurrentPeriodStart:  convertUnixToUTCTime(request.CurrentPeriodStart),
	})
	if err != nil {
		logger.Error("error on repository.create_membership", "request", request, "err", err)
		return err
	}

	return nil
}

func convertUnixToUTCTime(unixTime int64) time.Time {
	t := time.Unix(unixTime, 0)
	return t.UTC()
}
