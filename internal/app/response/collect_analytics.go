package response

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"net/http"
)

// CollectAnalyticsResponse is wrapper to proto file description to satisfy Renderer interface.
type CollectAnalyticsResponse struct {
	*desc.CollectAnalyticsResponse
}

func (resp *CollectAnalyticsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, int(resp.StatusCode))
	return nil
}