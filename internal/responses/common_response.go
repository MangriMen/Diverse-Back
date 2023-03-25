package responses

// BaseResponseBody represents the base response body for API.
// Contains error and message fields.
type BaseResponseBody struct {
	// required: true
	// example: false
	Error bool `json:"error"`

	// required: true
	Message *string `json:"message"`
}

// ErrorResponse represents the error response body for API.
// Contains error and message fields.
// swagger:response
type ErrorResponse struct {
	// in: body
	Body BaseResponseBody
}
