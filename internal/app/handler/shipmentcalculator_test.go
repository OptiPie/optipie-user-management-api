package handler

import (
	"context"
	"encoding/json"
	handlerresponse "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/response"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestShipmentCalculator(t *testing.T) {
	handlerName := "shipmentCalculatorHandler_test"
	lgr := zap.NewExample().Sugar()
	packSizes := []int{250, 500, 1000, 2000, 5000}
	tt := []struct {
		name                   string
		handler                *shipmentcalculator.ShipmentCalculator
		itemCount              string
		expectedItemPacks      map[int]int
		expectedHTTPStatusCode int
	}{
		{
			name: "happy path",
			handler: &shipmentcalculator.ShipmentCalculator{
				HandlerName: handlerName,
				Logger:      lgr,
				PackSizes:   packSizes,
			},
			itemCount:              "501",
			expectedItemPacks:      map[int]int{250: 1, 500: 1},
			expectedHTTPStatusCode: 200,
		},
		{
			name: "non integer item count as query parameter",
			handler: &shipmentcalculator.ShipmentCalculator{
				HandlerName: handlerName,
				Logger:      lgr,
				PackSizes:   packSizes,
			},
			itemCount:              "non-integer",
			expectedItemPacks:      nil,
			expectedHTTPStatusCode: 400,
		},
		{
			name: "negative item count as query parameter",
			handler: &shipmentcalculator.ShipmentCalculator{
				HandlerName: handlerName,
				Logger:      lgr,
				PackSizes:   packSizes,
			},
			itemCount:              "-50",
			expectedItemPacks:      nil,
			expectedHTTPStatusCode: 400,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/calculate/{itemCount}", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("itemCount", tc.itemCount)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := ShipmentCalculator(tc.handler)
			handler(w, r)

			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			defer response.Body.Close()

			statusCode := response.StatusCode
			if statusCode != tc.expectedHTTPStatusCode {
				t.Errorf("unexpected ShipmentCalculatorHandle() HTTP status code error = %v, wanted, %v", statusCode, tc.expectedHTTPStatusCode)
				return
			}

			shipmentCalculatorResponse := &handlerresponse.ShipmentCalculatorResponse{}
			json.Unmarshal(body, &shipmentCalculatorResponse)

			itemPacks := shipmentCalculatorResponse.ItemPacks
			if !reflect.DeepEqual(itemPacks, tc.expectedItemPacks) {
				t.Errorf("unexpected ShipmentCalculatorHandle() handle error = %v, wanted %v", itemPacks, tc.expectedItemPacks)
			}

		})
	}
}
