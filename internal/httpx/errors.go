package httpx

import (
	"encoding/json"
	"net/http"
)

type errorCode struct {
	value string
}

var (
	InvalidId        = errorCode{value: "invalid_id"}        // 400
	NotFound         = errorCode{value: "not_found"}         // 404
	InternalError    = errorCode{value: "internal_error"}    // 500
	MalformedJSON    = errorCode{value: "malformed_json"}    // 400
	BadRequest       = errorCode{value: "bad_request"}       // 400
	ValidationFailed = errorCode{value: "validation_failed"} // 422
	UnAuthenticated  = errorCode{value: "unauthenticated"}   // 401
	Forbidden        = errorCode{value: "forbidden"}         // 403
	Conflict         = errorCode{value: "conflict"}          // 409
	RateLimited      = errorCode{value: "rate_limited"}      // 429
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func Error(w http.ResponseWriter, status int, message string, code errorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			Message: message,
			Code:    code.value,
		},
	})

}
