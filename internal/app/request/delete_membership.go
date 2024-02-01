package request

import (
	"fmt"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"net/http"
)

type DeleteMembershipRequest struct {
	*desc.DeleteMembershipRequest
}

func (req DeleteMembershipRequest) Bind(r *http.Request) error {
	email := req.GetData().GetSupporterEmail()
	if email == "" {
		return fmt.Errorf("email is wrong")
	}
	return nil
}
