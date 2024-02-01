package usermanagementapi

import (
	apprequest "github.com/OptiPie/optipie-user-management-api/internal/app/request"
	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/usecase/handlers"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"net/http"
)

// UpdateMembership handles /api/v1/user/membership/update endpoint.
func (i *Implementation) UpdateMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := i.logger
	response := &appresponse.UpdateMembershipResponse{
		UpdateMembershipResponse: new(desc.UpdateMembershipResponse),
	}

	request := &apprequest.UpdateMembershipRequest{}

	if err := render.Bind(r, request); err != nil {
		logger.Error("error on binding request", "err", err)
		response.StatusCode = http.StatusBadRequest
		render.Render(w, r, response)
		return
	}

	data := request.GetData()

	err := i.updateMembershipHandler.HandleRequest(ctx, handlers.UpdateMembershipRequest{
		Paused:              data.GetPaused(),
		Status:              data.GetStatus(),
		Canceled:            data.GetCanceled(),
		Currency:            data.GetCurrency(),
		PspId:               data.GetPspId(),
		MembershipLevelId:   data.GetMembershipLevelId(),
		MembershipLevelName: data.GetMembershipLevelName(),
		StartedAt:           data.GetStartedAt(),
		CanceledAt:          data.GetCanceledAt(),
		NoteHidden:          data.GetNoteHidden(),
		SupportNote:         data.GetSupportNote(),
		SupporterName:       data.GetSupporterName(),
		SupporterEmail:      data.GetSupporterEmail(),
		CurrentPeriodEnd:    data.GetCurrentPeriodEnd(),
		SupporterFeedback:   data.GetSupporterFeedback(),
		CancelAtPeriodEnd:   data.GetCancelAtPeriodEnd(),
		CurrentPeriodStart:  data.GetCurrentPeriodStart(),
	})

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		render.Render(w, r, response)
		return
	}

	response.StatusCode = http.StatusOK
	render.Render(w, r, response)
}
