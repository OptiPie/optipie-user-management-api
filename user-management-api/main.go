package main

import (
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/app/prepare"
	usermanagementapi "github.com/OptiPie/optipie-user-management-api/internal/app/user-management-api"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	logger := prepare.SLogger()

	appConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("error on GetConfig, %v", err)
	}

	// init chi router
	r := chi.NewRouter()

	// prepare custom middlewares
	middlewares, err := prepare.Middlewares(appConfig)
	if err != nil {
		log.Fatalf("error on preparing middlewares %v", err)
	}

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Millisecond * time.Duration(appConfig.App.Timeout)))

	// construct handlers
	handlerCreateMembership, err := handlers.NewCreateMembership(handlers.NewCreateMembershipArgs{
		Logger: logger,
		Config: appConfig,
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
		r.Group(func(r chi.Router) {
			// add custom middlewares to handlers
			for _, mw := range middlewares {
				r.Use(mw)
			}
			r.Post("/user/membership", implementation.CreateMembership)
		})

	})

	log.Fatal(http.ListenAndServe(":3000", r))
}
