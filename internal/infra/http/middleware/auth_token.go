package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	appresponse "github.com/OptiPie/optipie-user-management-api/internal/app/response"
	"github.com/OptiPie/optipie-user-management-api/internal/infra/requestdata"
	desc "github.com/OptiPie/optipie-user-management-api/pkg/user-management-api"
	"github.com/go-chi/render"
)

type UserInfo struct {
	Email         string
	EmailVerified bool
}

type userInfoContextKey struct{}

var userInfoCtxKey = userInfoContextKey{}

// GoogleOAuthArgs contains arguments for Google OAuth authentication middleware.
type GoogleOAuthArgs struct {
	GoogleClientID string
}

// tokenInfoResponse represents the response from Google's tokeninfo endpoint
type tokenInfoResponse struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Scope         string `json:"scope"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"` // "true" or "false" as string
	ExpiresIn     string `json:"expires_in"`
	Exp           string `json:"exp"`
}

// GoogleOAuth validates Google OAuth access tokens from the Authorization header.
func GoogleOAuth(args GoogleOAuthArgs) Middleware {
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

			accessToken := parts[1]

			// Validate the access token with Google's tokeninfo endpoint
			tokenInfo, err := validateAccessToken(accessToken, args.GoogleClientID)
			if err != nil {
				slog.Error("invalid access token", "err", err)
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			// Extract email from the token info
			email := tokenInfo.Email
			if email == "" {
				slog.Error("email not found in token info")
				response.StatusCode = http.StatusUnauthorized
				render.Render(w, r, response)
				return
			}

			emailVerified, _ := strconv.ParseBool(tokenInfo.EmailVerified)

			userInfo := UserInfo{
				Email:         email,
				EmailVerified: emailVerified,
			}

			// Store user info in request context
			ctx := context.WithValue(r.Context(), userInfoCtxKey, &userInfo)
			r = r.WithContext(ctx)

			slog.Info("access token validated successfully", "email", userInfo.Email)

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

// validateAccessToken validates an OAuth access token with Google's tokeninfo endpoint
func validateAccessToken(accessToken, expectedClientID string) (*tokenInfoResponse, error) {
	// Call Google's tokeninfo endpoint
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	var tokenInfo tokenInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	// Verify the token was issued for our client ID
	if tokenInfo.Aud != expectedClientID {
		return nil, fmt.Errorf("token audience mismatch: expected %s, got %s", expectedClientID, tokenInfo.Aud)
	}

	return &tokenInfo, nil
}
