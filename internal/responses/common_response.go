package responses

// BaseResponseBody is base body for application response.
// Includes the error and message fields.
type BaseResponseBody struct {
	// required: true
	// example: false
	Error bool `json:"error" validate:"required"`

	// required: true
	Message interface{} `json:"message" validate:"required"`
}

type errorResponseBody struct {
	BaseResponseBody

	// required: true
	// example: true
	Error bool `json:"error" validate:"required"`
}

// ErrorResponse contains the error data.
// swagger:response
type ErrorResponse struct {
	// in: body
	Body errorResponseBody
}

// SuccessResponse represents success response for request.
// swagger:response
type SuccessResponse string
