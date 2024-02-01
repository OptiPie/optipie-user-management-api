package request

import (
	"fmt"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type UpdateMembershipRequest struct {
	*desc.UpdateMembershipRequest
}

func (req UpdateMembershipRequest) Bind(r *http.Request) error {
	email := req.GetData().GetSupporterEmail()
	if email == "" {
		return fmt.Errorf("email is wrong")
	}
	return nil
}
