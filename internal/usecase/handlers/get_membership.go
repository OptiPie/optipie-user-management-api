package handlers

import (
	"context"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"log/slog"
	"strconv"
	"time"
)

const (
	membershipStatusActive = "active"
)

// GetMembershipHandler is an abstraction for GetMembership use-case handler.
type GetMembershipHandler interface {
	HandleRequest(ctx context.Context, request GetMembershipRequest) (*GetMembershipResponse, error)
}

type NewGetMembershipArgs struct {
	Logger     *slog.Logger
	Config     *config.Config
	Repository domain.Repository
}

func NewGetMembership(args NewGetMembershipArgs) (*GetMembership, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if args.Repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	return &GetMembership{
		logger:     args.Logger,
		config:     args.Config,
		repository: args.Repository,
	}, nil
}

// GetMembership is a request handler with all dependencies initialized.
type GetMembership struct {
	logger     *slog.Logger
	config     *config.Config
	repository domain.Repository
}

// GetMembershipRequest represents necessary GET /api/v1/user/membership/{email} request data for handler.
type GetMembershipRequest struct {
	Email string
}

type GetMembershipResponse struct {
	Email                      string
	IsMembershipExists         bool
	IsMembershipActive         bool
	IsMembershipPaused         *bool
	IsMembershipCanceled       *bool
	CurrentMembershipPeriodEnd *int64
}

func (h *GetMembership) HandleRequest(ctx context.Context, request GetMembershipRequest) (*GetMembershipResponse, error) {
	logger := h.logger
	repository := h.repository

	membership, err := repository.GetMembershipByEmail(ctx, request.Email)

	if err != nil {
		logger.Error("error on repository.get_membership_by_email", "request", request, "err", err)
		return nil, err
	}

	// meaning membership entity for email doesn't exist
	if membership.Id == 0 {
		return &GetMembershipResponse{
			Email:              membership.SupporterEmail,
			IsMembershipExists: false,
			IsMembershipActive: false,
		}, nil
	}

	paused, _ := strconv.ParseBool(membership.Paused)
	canceled, _ := strconv.ParseBool(membership.Canceled)

	// check if membership is valid
	if time.Now().After(membership.CurrentPeriodEnd) || membership.Status != membershipStatusActive {
		logger.Warn("membership is not active", "membership", membership)
		return &GetMembershipResponse{
			Email:                membership.SupporterEmail,
			IsMembershipExists:   true,
			IsMembershipActive:   false,
			IsMembershipPaused:   &paused,
			IsMembershipCanceled: &canceled,
		}, nil
	}

	currentPeriodEndTimestamp := membership.CurrentPeriodEnd.Unix()

	response := &GetMembershipResponse{
		Email:                      membership.SupporterEmail,
		IsMembershipExists:         true,
		IsMembershipActive:         true,
		IsMembershipPaused:         &paused,
		IsMembershipCanceled:       &canceled,
		CurrentMembershipPeriodEnd: &currentPeriodEndTimestamp,
	}

	return response, nil
}
