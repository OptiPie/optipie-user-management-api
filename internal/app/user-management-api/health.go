package usermanagementapi

import (
	"net/http"
)

const HealthEndpoint = "/health"

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
