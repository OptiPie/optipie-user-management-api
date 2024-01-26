package request

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type UpdateMembershipRequest struct {
	*desc.UpdateMembershipRequest
}

func (req UpdateMembershipRequest) Bind(r *http.Request) error {
	return nil
}
