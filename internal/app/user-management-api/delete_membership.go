package usermanagementapi

import (
	"github.com/OptiPie/optipie-user-management-api/internal/app/cerrors"
	apprequest "github.com/OptiPie/optipie-user-management-api/internal/app/request"
	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"net/http"
)

// DeleteMembership handles /api/v1/user/membership/delete endpoint.
func (i *Implementation) DeleteMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := i.logger
	response := &appresponse.DeleteMembershipResponse{
		DeleteMembershipResponse: new(desc.DeleteMembershipResponse),
	}

	request := &apprequest.DeleteMembershipRequest{}

	if err := render.Bind(r, request); err != nil {
		logger.Error("error on binding request", "err", err)
		response.StatusCode = http.StatusBadRequest
		render.Render(w, r, response)
		return
	}

	data := request.GetData()

	err := i.deleteMembershipHandler.HandleRequest(ctx, handlers.DeleteMembershipRequest{
		Email: data.SupporterEmail,
	})

	// custom error handling in case entity for given email is not found
	if customErr, ok := err.(*cerrors.CustomError); ok {
		if customErr.TypesMap[cerrors.ConditionalCheckFailedException] {
			response.StatusCode = http.StatusNotFound
			render.Render(w, r, response)
			return
		}
	}

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		render.Render(w, r, response)
		return
	}

	response.StatusCode = http.StatusOK
	render.Render(w, r, response)
}
