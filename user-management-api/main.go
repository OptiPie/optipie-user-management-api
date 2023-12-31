package main

import (
	"context"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/app/prepare"
	usermanagementapi "github.com/OptiPie/optipie-user-management-api/internal/app/user-management-api"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := prepare.SLogger()

	appConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("error on GetConfig, %v", err)
	}

	// init chi router
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Millisecond * time.Duration(appConfig.App.Timeout)))
	r.Use(MyMiddleware)

	// construct handlers
	handlerCreateMembership, err := handlers.NewCreateMembership(handlers.NewCreateMembershipArgs{
		Logger:  logger,
		Config:  appConfig,
		AppName: "test",
	})

	if err != nil {
		log.Fatalf("error on NewCreateMembership, %v", err)
	}

	implementation, err := usermanagementapi.NewUserManagementAPI(
		usermanagementapi.NewUserManagementAPIArgs{
			Logger:                  logger,
			Config:                  appConfig,
			CreateMembershipHandler: handlerCreateMembership,
		})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get(usermanagementapi.HealthEndpoint, usermanagementapi.Health)
		r.Post("/user/membership", implementation.CreateMembership)

	})

	log.Fatal(http.ListenAndServe(":3000", r))
}

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), "user", "123")

		// call the next user-management-api in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
