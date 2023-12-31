package response

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"net/http"
)

// CreateMembershipResponse is wrapper to proto file description to satisfy Renderer interface.
type CreateMembershipResponse struct {
	*desc.CreateMembershipResponse
}

func (e *CreateMembershipResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, int(e.StatusCode))
	return nil
}
