package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/infra/requestdata"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

// Middleware is an explicit middleware type for better readability.
type Middleware func(next http.Handler) http.Handler

// AuthArgs is for passing all args in human-readable format.
type AuthArgs struct {
	SecretKey string
}

func Auth(args AuthArgs) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := &appresponse.CreateMembershipResponse{
				CreateMembershipResponse: new(desc.CreateMembershipResponse),
			}

			secretKey := args.SecretKey
			requestHeaders := requestdata.GetRequestHeaders(r)
			requestBodyMAC := []byte(requestHeaders.BMCSignature)

			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("error on reading request body", "err", err)
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			mac := hmac.New(sha256.New, []byte(secretKey))
			mac.Write(requestBody)
			expectedMAC := mac.Sum(nil)

			isMACValid := hmac.Equal(requestBodyMAC, expectedMAC)
			
			slog.Info("isMACValid", "", isMACValid)
			slog.Info("auth middleware compare request body signatures", "expectedMAC", expectedMAC, "requestBodyMAC", requestBodyMAC)

			next.ServeHTTP(w, r)
		})
	}
}
