package prepare

import (
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	mw "github.com/OptiPie/optipie-user-management-api/internal/infra/http/middleware"
)

// Middlewares prepares all the necessary mws for the app.
func Middlewares(config *config.Config) ([]mw.Middleware, error) {
	var middlewares []mw.Middleware
	authMiddleware := mw.Auth(mw.AuthArgs{
		MembershipStartedSecretKey:   config.App.WebHookKeys.MembershipStarted,
		MembershipUpdatedSecretKey:   config.App.WebHookKeys.MembershipUpdated,
		MembershipCancelledSecretKey: config.App.WebHookKeys.MembershipCancelled,
	})

	middlewares = append(middlewares, authMiddleware)
	return middlewares, nil
}
