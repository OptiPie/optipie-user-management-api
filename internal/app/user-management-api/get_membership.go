package usermanagementapi

import (
	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

const (
	maxEmailCharacters = 320
)

func (i *Implementation) GetMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := i.logger
	response := &appresponse.GetMembershipResponse{
		GetMembershipResponse: new(desc.GetMembershipResponse),
	}

	var email string

	if email = chi.URLParam(r, "email"); email == "" || len(email) > maxEmailCharacters {
		logger.Error("email is wrong", "email", email)
		response.StatusCode = http.StatusBadRequest
		render.Render(w, r, response)
		return
	}

	membershipResponse, err := i.getMembershipHandler.HandleRequest(ctx, handlers.GetMembershipRequest{
		Email: email,
	})

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		render.Render(w, r, response)
		return
	}

	if !membershipResponse.IsMembershipExists {
		response.StatusCode = http.StatusNotFound
		render.Render(w, r, response)
		return
	}

	responseData := &desc.GetMembershipResponse_Data{
		Email:                      &membershipResponse.Email,
		IsMembershipActive:         &membershipResponse.IsMembershipActive,
		IsMembershipPaused:         membershipResponse.IsMembershipPaused,
		IsMembershipCanceled:       membershipResponse.IsMembershipCanceled,
		CurrentMembershipPeriodEnd: membershipResponse.CurrentMembershipPeriodEnd,
	}

	response.Data = responseData
	response.StatusCode = http.StatusOK
	render.Render(w, r, response)
}
