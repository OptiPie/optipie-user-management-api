package requestdata

import (
	"log/slog"
	"net/http"
)

const (
	headerBMCSignature      = "X-BMC-Signature"
	headerRequestID         = "X-Request-ID"
	headerOriginalUserAgent = "X-Original-User-Agent"
	headerUserAgent         = "User-Agent"
)

// RequestHeaders stores http request headers.
type RequestHeaders struct {
	BMCSignature      string
	RequestID         string
	OriginalUserAgent string
	UserAgent         string
}

// GetRequestHeaders is for getting request data from http request headers.
func GetRequestHeaders(req *http.Request) RequestHeaders {
	if req == nil {
		return RequestHeaders{}
	}

	slog.Info("All headers are : ", "RequestHeaders", req.Header)

	requestID := req.Header.Get(headerRequestID)

	return RequestHeaders{
		BMCSignature:      req.Header.Get(headerBMCSignature),
		RequestID:         requestID,
		OriginalUserAgent: req.Header.Get(headerOriginalUserAgent),
		UserAgent:         req.UserAgent(),
	}
}
