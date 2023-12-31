package request

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type CreateMembershipRequest struct {
	*desc.CreateMembershipRequest
}

func (req CreateMembershipRequest) Bind(r *http.Request) error {

	return nil
}
