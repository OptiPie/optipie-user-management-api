package request

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type GetMembershipRequest struct {
	*desc.GetMembershipRequest
}

func (req GetMembershipRequest) Bind(r *http.Request) error {
	return nil
}
