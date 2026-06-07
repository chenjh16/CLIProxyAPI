package executor

import (
	"net/http"
	"testing"
)

func TestStatusErrTreatsRateLimitCooldownAsQuota(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "error code",
			body: `{
				"error": {
					"message": "cooldown",
					"type": "invalid_request_error",
					"code": "rate_limit_cooldown"
				},
				"message": "cooldown",
				"code": "rate_limit_cooldown",
				"limit_type": "cooldown"
			}`,
		},
		{
			name: "top-level code fallback",
			body: `{
				"error": {
					"message": "cooldown",
					"type": "invalid_request_error"
				},
				"message": "cooldown",
				"code": "rate_limit_cooldown",
				"limit_type": "cooldown"
			}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := statusErr{code: http.StatusBadRequest, msg: tc.body}
			if got := err.StatusCode(); got != http.StatusTooManyRequests {
				t.Fatalf("StatusCode() = %d, want %d", got, http.StatusTooManyRequests)
			}
		})
	}
}

func TestStatusErrRateLimitCooldownRequiresStructuredMatch(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing error object",
			body: `{"message":"invalid_request_error rate_limit_cooldown","code":"rate_limit_cooldown"}`,
		},
		{
			name: "wrong error type",
			body: `{"error":{"type":"server_error","code":"rate_limit_cooldown"}}`,
		},
		{
			name: "wrong error code",
			body: `{"error":{"type":"invalid_request_error","code":"context_length_exceeded"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := statusErr{code: http.StatusBadRequest, msg: tc.body}
			if got := err.StatusCode(); got != http.StatusBadRequest {
				t.Fatalf("StatusCode() = %d, want %d", got, http.StatusBadRequest)
			}
		})
	}
}
