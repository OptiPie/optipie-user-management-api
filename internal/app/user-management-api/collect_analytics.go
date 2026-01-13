package usermanagementapi

import (
	apprequest "github.com/OptiPie/optipie-user-management-api/internal/app/request"
	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"net/http"
)

// CollectAnalytics handles /api/v1/analytics/collect endpoint.
func (i *Implementation) CollectAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := i.logger
	response := &appresponse.CollectAnalyticsResponse{
		CollectAnalyticsResponse: new(desc.CollectAnalyticsResponse),
	}

	request := &apprequest.CollectAnalyticsRequest{}

	if err := render.Bind(r, request); err != nil {
		logger.Error("error on binding request", "err", err)
		response.StatusCode = http.StatusBadRequest
		render.Render(w, r, response)
		return
	}

	err := i.collectAnalyticsHandler.HandleRequest(ctx, handlers.CollectAnalyticsRequest{
		StrategyName:   request.GetStrategyName(),
		StrategySymbol: request.GetStrategySymbol(),
		StrategyPeriod: request.GetStrategyPeriod(),
	})

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		render.Render(w, r, response)
		return
	}

	response.StatusCode = http.StatusOK
	render.Render(w, r, response)
}