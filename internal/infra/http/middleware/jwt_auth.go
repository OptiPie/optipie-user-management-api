package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/infra/requestdata"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
	"google.golang.org/api/idtoken"
)

type UserInfo struct {
	Email         string
	EmailVerified bool
}

type userInfoContextKey struct{}

var userInfoCtxKey = userInfoContextKey{}

// JWTAuthArgs contains arguments for JWT authentication middleware.
type JWTAuthArgs struct {
	GoogleClientID string
}

// JWTAuth validates Google OAuth JWT tokens from the Authorization header.
func JWTAuth(args JWTAuthArgs) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := &appresponse.GetMembershipResponse{
				GetMembershipResponse: new(desc.GetMembershipResponse),
			}

			requestHeaders := requestdata.GetRequestHeaders(r)
			authHeader := requestHeaders.Authorization

			if authHeader == "" {
				slog.Error("missing authorization header")
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			// Extract token from "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				slog.Error("invalid authorization header format")
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			idToken := parts[1]

			// Validate the token with Google
			ctx := context.Background()
			payload, err := idtoken.Validate(ctx, idToken, args.GoogleClientID)
			if err != nil {
				slog.Error("invalid JWT token", "err", err)
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			// Extract email from the token payload
			email, ok := payload.Claims["email"].(string)
			if !ok {
				slog.Error("email claim not found in token")
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			emailVerified, _ := payload.Claims["email_verified"].(bool)

			userInfo := UserInfo{
				Email:         email,
				EmailVerified: emailVerified,
			}

			// Store user info in request context
			ctx = context.WithValue(r.Context(), userInfoCtxKey, &userInfo)
			r = r.WithContext(ctx)

			slog.Info("JWT token validated successfully", "email", userInfo.Email)

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserInfoFromContext retrieves user information from the request context.
func GetUserInfoFromContext(ctx context.Context) (*UserInfo, bool) {
	userInfo := ctx.Value(userInfoCtxKey)
	if userInfo == nil {
		return nil, false
	}
	if info, ok := userInfo.(*UserInfo); ok {
		return info, true
	}
	return nil, false
}