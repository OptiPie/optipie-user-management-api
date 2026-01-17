package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/domain"
)

// CollectAnalyticsHandler is an abstraction for analytics collection use-case handler.
type CollectAnalyticsHandler interface {
	HandleRequest(ctx context.Context, request CollectAnalyticsRequest) error
}

type NewCollectAnalyticsArgs struct {
	Logger     *slog.Logger
	Config     *config.Config
	Repository domain.Repository
}

func NewCollectAnalytics(args NewCollectAnalyticsArgs) (*CollectAnalytics, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if args.Repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	return &CollectAnalytics{
		logger:     args.Logger,
		config:     args.Config,
		repository: args.Repository,
	}, nil
}

// CollectAnalytics is a request handler with all dependencies initialized.
type CollectAnalytics struct {
	logger     *slog.Logger
	config     *config.Config
	repository domain.Repository
}

// CollectAnalyticsRequest represents necessary POST /api/v1/analytics/collect request data for handler.
type CollectAnalyticsRequest struct {
	StrategyName      string
	StrategySymbol    string
	StrategyPeriod    string
	StrategyDateRange string
}

func (h *CollectAnalytics) HandleRequest(ctx context.Context, request CollectAnalyticsRequest) error {
	logger := h.logger
	repository := h.repository

	timestamp := time.Now().UTC().UnixNano()

	err := repository.CreateAnalytics(ctx, domain.CreateAnalyticsArgs{
		Timestamp:         timestamp,
		StrategyName:      request.StrategyName,
		StrategySymbol:    request.StrategySymbol,
		StrategyPeriod:    request.StrategyPeriod,
		StrategyDateRange: request.StrategyDateRange,
	})
	if err != nil {
		logger.Error("error on repository.create_analytics", "request", request, "err", err)
		return err
	}

	return nil
}
