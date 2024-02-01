package handlers

import (
	"context"
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
	"log/slog"
)

// DeleteMembershipHandler is an abstraction for MembershipCancelled use-case handler.
type DeleteMembershipHandler interface {
	HandleRequest(ctx context.Context, request DeleteMembershipRequest) error
}

type NewDeleteMembershipArgs struct {
	Logger     *slog.Logger
	Config     *config.Config
	Repository domain.Repository
}

func NewDeleteMembership(args NewDeleteMembershipArgs) (*DeleteMembership, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if args.Repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	return &DeleteMembership{
		logger:     args.Logger,
		config:     args.Config,
		repository: args.Repository,
	}, nil
}

// DeleteMembership is a request handler with all dependencies initialized.
type DeleteMembership struct {
	logger     *slog.Logger
	config     *config.Config
	repository domain.Repository
}

// DeleteMembershipRequest represents necessary GET /api/v1/user/membership/delete request data for handler.
type DeleteMembershipRequest struct {
	Email string
}

func (h *DeleteMembership) HandleRequest(ctx context.Context, request DeleteMembershipRequest) error {
	logger := h.logger
	repository := h.repository

	err := repository.DeleteMembershipByEmail(ctx, request.Email)

	if err != nil {
		logger.Error("error on repository.delete_membership_by_email", "request", request, "err", err)
		return err
	}

	return nil
}
