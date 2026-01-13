package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OptiPie/optipie-user-management-api/internal/app/config"
	"github.com/OptiPie/optipie-user-management-api/internal/app/prepare"
	usermanagementapi "github.com/OptiPie/optipie-user-management-api/internal/app/user-management-api"
	dynamorepo "github.com/OptiPie/optipie-user-management-api/internal/infra/dynamodb"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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

	googleOAuthMiddleware := prepare.GoogleOAuthMiddleware(appConfig)

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
	repository := dynamorepo.NewRepository(svc, appConfig.Aws.Dynamodb.Membership.TableName, appConfig.Aws.Dynamodb.Analytics.TableName)

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Millisecond * time.Duration(appConfig.App.Timeout)))

	// handlers
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
	if err != nil {
		log.Fatalf("error on NewGetMembership, %v", err)
	}

	handlerUpdateMembership, err := handlers.NewUpdateMembership(handlers.NewUpdateMembershipArgs{
		Logger:     logger,
		Config:     appConfig,
		Repository: repository,
	})
	if err != nil {
		log.Fatalf("error on NewUpdateMembership, %v", err)
	}

	handlerDeleteMembership, err := handlers.NewDeleteMembership(handlers.NewDeleteMembershipArgs{
		Logger:     logger,
		Config:     appConfig,
		Repository: repository,
	})
	if err != nil {
		log.Fatalf("error on NewDeleteMembership, %v", err)
	}

	handlerCollectAnalytics, err := handlers.NewCollectAnalytics(handlers.NewCollectAnalyticsArgs{
		Logger:     logger,
		Config:     appConfig,
		Repository: repository,
	})
	if err != nil {
		log.Fatalf("error on NewCollectAnalytics, %v", err)
	}

	implementation, err := usermanagementapi.NewUserManagementAPI(
		usermanagementapi.NewUserManagementAPIArgs{
			Logger:                  logger,
			Config:                  appConfig,
			CreateMembershipHandler: handlerCreateMembership,
			GetMembershipHandler:    handlerGetMembership,
			UpdateMembershipHandler: handlerUpdateMembership,
			DeleteMembershipHandler: handlerDeleteMembership,
			CollectAnalyticsHandler: handlerCollectAnalytics,
		})

	if err != nil {
		log.Fatalf("error on NewUserManagementAPI, %v", err)
	}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get(usermanagementapi.HealthEndpoint, usermanagementapi.Health)
		r.Route("/user/membership", func(r chi.Router) {
			// rate limit by IP, 100 requests per minute
			r.Use(httprate.LimitByIP(100, 1*time.Minute))
			r.Group(func(r chi.Router) {
				// add custom middlewares to webhook handlers
				for _, mw := range middlewares {
					r.Use(mw)
				}
				r.Post("/create", implementation.CreateMembership)
				r.Post("/update", implementation.UpdateMembership)
				r.Post("/delete", implementation.DeleteMembership)
			})
			// GET doesn't require any auth or custom mw logic
			r.Route("/{email}", func(r chi.Router) {
				r.Get("/", implementation.GetMembership)
			})
		})
		r.Route("/analytics", func(r chi.Router) {
			// rate limit by IP, 100 requests per minute
			r.Use(httprate.LimitByIP(100, 1*time.Minute))
			r.Use(googleOAuthMiddleware)
			r.Post("/collect", implementation.CollectAnalytics)
		})
	})

	log.Fatal(http.ListenAndServe(":3000", r))
}
