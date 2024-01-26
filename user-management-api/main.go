package main

import (
	"context"
	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/app/prepare"
	usermanagementapi "github.com/OptiPie/optipie-user-management-api/internal/app/user-management-api"
	dynamorepo "github.com/OptiPie/optipie-user-management-api/internal/infra/dynamodb"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	ctx := context.Background()
	// init chi router
	r := chi.NewRouter()

	// prepare custom middlewares
	middlewares, err := prepare.Middlewares(appConfig)
	if err != nil {
		log.Fatalf("error on preparing middlewares %v", err)
	}

	var awsCfg aws.Config
	if appConfig.App.IsLocalDevelopment {
		logger.Warn("Warning, local development environment!")
		awsCfg, err = prepare.LocalAwsConfig(ctx)
	} else {
		awsCfg, err = prepare.AwsConfig(ctx)
	}

	if err != nil {
		log.Fatalf("prepare aws config error: %v", err)
	}

	svc := prepare.Dynamodb(awsCfg)
	repository := dynamorepo.NewRepository(svc)

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Millisecond * time.Duration(appConfig.App.Timeout)))

	// construct handlers
	handlerCreateMembership, err := handlers.NewCreateMembership(handlers.NewCreateMembershipArgs{
		Logger:     logger,
		Config:     appConfig,
		Repository: repository,
	})
	if err != nil {
		log.Fatalf("error on NewCreateMembership, %v", err)
	}

	handlerGetMembership, err := handlers.NewGetMembership(handlers.NewGetMembershipArgs{
		Logger:     logger,
		Config:     appConfig,
		Repository: repository,
	})

	handlerUpdateMembership, err := handlers.NewUpdateMembership(handlers.NewUpdateMembershipArgs{
		Logger:     logger,
		Config:     appConfig,
		Repository: repository,
	})

	if err != nil {
		log.Fatalf("error on NewGetMembership, %v", err)
	}

	implementation, err := usermanagementapi.NewUserManagementAPI(
		usermanagementapi.NewUserManagementAPIArgs{
			Logger:                  logger,
			Config:                  appConfig,
			CreateMembershipHandler: handlerCreateMembership,
			GetMembershipHandler:    handlerGetMembership,
			UpdateMembershipHandler: handlerUpdateMembership,
		})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get(usermanagementapi.HealthEndpoint, usermanagementapi.Health)
		r.Group(func(r chi.Router) {
			// add custom middlewares to handlers
			for _, mw := range middlewares {
				r.Use(mw)
			}
			r.Route("/user/membership", func(r chi.Router) {
				r.Post("/create", implementation.CreateMembership)
				r.Post("/update", implementation.UpdateMembership)
				r.Post("/delete", nil)
				r.Route("/{email}", func(r chi.Router) {
					r.Get("/", implementation.GetMembership)
				})
			})

		})
	})

	log.Fatal(http.ListenAndServe(":3000", r))
}
