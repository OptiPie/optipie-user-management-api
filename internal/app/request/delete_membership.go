package request

import (
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type DeleteMembershipRequest struct {
	*desc.DeleteMembershipRequest
}

func (req DeleteMembershipRequest) Bind(r *http.Request) error {
	return nil
}
