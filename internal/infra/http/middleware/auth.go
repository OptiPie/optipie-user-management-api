package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/infra/requestdata"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

const (
	membershipTypeStarted   = "membership.started"
	membershipTypeUpdated   = "membership.updated"
	membershipTypeCancelled = "membership.cancelled"
)

// Middleware is an explicit middleware type for better readability.
type Middleware func(next http.Handler) http.Handler

// AuthArgs is for passing all args in human-readable format.
type AuthArgs struct {
	MembershipStartedSecretKey   string
	MembershipUpdatedSecretKey   string
	MembershipCancelledSecretKey string
}

type MetaData struct {
	Type     string `json:"type"`
	LiveMode bool   `json:"live_mode"`
	Attempt  int    `json:"attempt"`
}

func Auth(args AuthArgs) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := &appresponse.CreateMembershipResponse{
				CreateMembershipResponse: new(desc.CreateMembershipResponse),
			}

			secretKey := ""
			requestHeaders := requestdata.GetRequestHeaders(r)
			requestBodyMAC := []byte(requestHeaders.SignatureSha256)

			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("error on reading request body", "err", err)
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			slog.Info("exact request body is: ", "request_body", requestBody)

			metaData := &MetaData{}
			err = json.Unmarshal(requestBody, metaData)
			if err != nil {
				slog.Error("error on unmarshalling request body", "err", err)
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			switch metaData.Type {
			case membershipTypeStarted:
				secretKey = args.MembershipStartedSecretKey
			case membershipTypeUpdated:
				secretKey = args.MembershipUpdatedSecretKey
			case membershipTypeCancelled:
				secretKey = args.MembershipCancelledSecretKey
			}

			mac := hmac.New(sha256.New, []byte(secretKey))
			mac.Write(requestBody)
			expectedMAC := mac.Sum(nil)

			hmacHex := hex.EncodeToString(expectedMAC)
			fmt.Printf("%v,%v", hmacHex)

			isMACValid := hmac.Equal(requestBodyMAC, []byte(hmacHex))

			slog.Info("isMACValid", "", isMACValid)
			slog.Info("auth middleware compare request body signatures", "expectedMAC", expectedMAC, "requestBodyMAC", requestBodyMAC)

			next.ServeHTTP(w, r)
		})
	}
}
