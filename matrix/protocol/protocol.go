package protocol

import m "github.com/andir/matrix_server/matrix"

type ErrorResponse struct {
	Errcode m.MErrorCode `json:"errcode"`
	Error string `json:"error"`
}

type RateLimitResponse struct {
	ErrorResponse
	RetryAfterMs uint `json:"retry_after_ms,omitempty"`
}
