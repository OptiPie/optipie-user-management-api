package requestdata

import (
	"log/slog"
	"net/http"
)

const (
	headerSignatureSha256   = "X-Signature-Sha256"
	headerRequestID         = "X-Request-ID"
	headerOriginalUserAgent = "X-Original-User-Agent"
	headerUserAgent         = "User-Agent"
)

// RequestHeaders stores http request headers.
type RequestHeaders struct {
	SignatureSha256   string
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
		SignatureSha256:   req.Header.Get(headerSignatureSha256),
		RequestID:         requestID,
		OriginalUserAgent: req.Header.Get(headerOriginalUserAgent),
		UserAgent:         req.UserAgent(),
	}
}
