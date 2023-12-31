package usermanagementapi

import (
	"fmt"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	"log/slog"
)

type Implementation struct {
	logger                  *slog.Logger
	config                  *config.Config
	createMembershipHandler handlers.CreateMembershipHandler
}

func NewUserManagementAPI(args NewUserManagementAPIArgs) (*Implementation, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if args.CreateMembershipHandler == nil {
		return nil, fmt.Errorf("membershipStartedHandler is required")
	}
	return &Implementation{
		logger:                  args.Logger,
		config:                  args.Config,
		createMembershipHandler: args.CreateMembershipHandler,
	}, nil
}

type NewUserManagementAPIArgs struct {
	Logger                  *slog.Logger
	Config                  *config.Config
	CreateMembershipHandler handlers.CreateMembershipHandler
}
