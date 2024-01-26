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
	response := &appresponse.CreateMembershipResponse{
		CreateMembershipResponse: new(desc.CreateMembershipResponse),
	}

	request := &apprequest.CreateMembershipRequest{}

	if err := render.Bind(r, request); err != nil {
		logger.Error("Error, %v", err)
		response.StatusCode = http.StatusBadRequest
		render.Render(w, r, response)
		return
	}
	logger.Info("Request is here pal!", request)

	data := request.GetData()

	if data.SupporterEmail == "" {
		logger.Error("supporter email can't be nil", "request", request)
		response.StatusCode = http.StatusBadRequest
		render.Render(w, r, response)
		return
	}

	err := i.createMembershipHandler.HandleRequest(ctx, handlers.CreateMemberShipRequest{
		Type:                request.GetType(),
		LiveMode:            request.GetLiveMode(),
		Attempt:             request.GetAttempt(),
		Created:             request.GetCreated(),
		EventId:             request.GetEventId(),
		Id:                  data.GetId(),
		Amount:              data.GetAmount(),
		Object:              data.GetObject(),
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
		SupporterId:         data.GetSupporterId(),
		SupporterEmail:      data.GetSupporterEmail(),
		CurrentPeriodEnd:    data.GetCurrentPeriodEnd(),
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
