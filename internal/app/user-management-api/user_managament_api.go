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
	getMembershipHandler    handlers.GetMembershipHandler
	updateMembershipHandler handlers.UpdateMembershipHandler
	deleteMembershipHandler handlers.DeleteMembershipHandler
}

func NewUserManagementAPI(args NewUserManagementAPIArgs) (*Implementation, error) {
	if args.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if args.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if args.CreateMembershipHandler == nil {
		return nil, fmt.Errorf("createMembershipHandler is required")
	}
	if args.GetMembershipHandler == nil {
		return nil, fmt.Errorf("getMembershipHandler is required")
	}
	if args.UpdateMembershipHandler == nil {
		return nil, fmt.Errorf("updateMembershipHandler is required")
	}
	if args.DeleteMembershipHandler == nil {
		return nil, fmt.Errorf("deleteMembershipHandler is required")
	}
	return &Implementation{
		logger:                  args.Logger,
		config:                  args.Config,
		createMembershipHandler: args.CreateMembershipHandler,
		getMembershipHandler:    args.GetMembershipHandler,
		updateMembershipHandler: args.UpdateMembershipHandler,
		deleteMembershipHandler: args.DeleteMembershipHandler,
	}, nil
}

type NewUserManagementAPIArgs struct {
	Logger                  *slog.Logger
	Config                  *config.Config
	CreateMembershipHandler handlers.CreateMembershipHandler
	GetMembershipHandler    handlers.GetMembershipHandler
	UpdateMembershipHandler handlers.UpdateMembershipHandler
	DeleteMembershipHandler handlers.DeleteMembershipHandler
}
