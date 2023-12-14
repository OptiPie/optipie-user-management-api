package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	logger := prepare.ZapLogger()

	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Fatalf("error on GetConfig, %v", err)
	}

	// init chi router
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Millisecond * time.Duration(appConfig.App.Timeout)))
	r.Use(MyMiddleware)

	// business handlers
	handler, err := prepare.NewShipmentCalculatorHandler("shipmentCalculatorHandler", logger, appConfig.App.PackSizes)
	if err != nil {
		logger.Fatalf("error on prepare shipmentCalculatorHandler, %v", err)
	}

	// app handlers
	r.Get(apphandler.HealthEndpoint, apphandler.Health)

	shipmentCalculatorRoute := fmt.Sprintf("%v/{%v}",
		apphandler.ShipmentCalculatorEndpoint, apphandler.ShipmentCalculatorURLParam)

	r.Get(shipmentCalculatorRoute, apphandler.ShipmentCalculator(handler))

	// swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(":3000/swagger/doc.json")), //The url pointing to API definition
	)

	logger.Fatal(http.ListenAndServe(":3000", r))
}

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  // create new context from `r` request context, and assign key `"user"`
	  // to value of `"123"`
	  ctx := context.WithValue(r.Context(), "user", "123")
  
	  // call the next handler in the chain, passing the response writer and
	  // the updated request object with the new context value.
	  //
	  // note: context.Context values are nested, so any previously set
	  // values will be accessible as well, and the new `"user"` key
	  // will be accessible from this point forward.
	  next.ServeHTTP(w, r.WithContext(ctx))
	})
  }
