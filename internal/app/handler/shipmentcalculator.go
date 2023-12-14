package handler

import (
	"fmt"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/response"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

const ShipmentCalculatorURLParam = "itemCount"
const ShipmentCalculatorEndpoint = "/calculate"

func ShipmentCalculator(handler shipmentcalculator.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		itemCountStr := chi.URLParam(r, ShipmentCalculatorURLParam)
		itemCount, err := strconv.Atoi(itemCountStr)

		if err != nil {
			render.Render(w, r, response.ErrInvalidRequest(err))
			return
		}

		if itemCount <= 0 {
			err = fmt.Errorf("item count should be greater than 0")
			render.Render(w, r, response.ErrInvalidRequest(err))
			return
		}

		packCount, err := handler.Handle(ctx, itemCount)

		if err != nil {
			render.Render(w, r, response.ErrInternalServer(err))
			return
		}

		shipmentCalculatorResponse := response.NewShipmentCalculatorResponse(http.StatusOK, "Number of Item Packs", packCount)

		err = render.Render(w, r, shipmentCalculatorResponse)
		if err != nil {
			render.Render(w, r, response.ErrInternalServer(err))
		}
	}
}
