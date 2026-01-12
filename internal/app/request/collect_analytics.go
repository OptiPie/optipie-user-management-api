package request

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type CollectAnalyticsRequest struct {
	*desc.CollectAnalyticsRequest
}

func (req CollectAnalyticsRequest) Bind(r *http.Request) error {
	return nil
}