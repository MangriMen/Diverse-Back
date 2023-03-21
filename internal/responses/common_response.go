package responses

type BaseResponseBody struct {
	// required: true
	// example: false
	Error bool `json:"error"`

	// required: true
	Message *string `json:"message"`
}

// swagger:response
type ErrorResponse struct {
	// in: body
	Body BaseResponseBody
}
