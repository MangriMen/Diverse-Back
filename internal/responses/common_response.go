package responses

// BaseResponseBody represents the base response body for API.
// Contains error and message fields.
type BaseResponseBody struct {
	// required: true
	// example: false
	Error bool `json:"error"`

	// required: true
	Message interface{} `json:"message"`
}

// ErrorResponse represents the error response body for API.
// Contains error and message fields.
// swagger:response
type ErrorResponse struct {
	// in: body
	Body BaseResponseBody
}

// SuccessResponse represents success response for request.
// swagger:response
type SuccessResponse string
