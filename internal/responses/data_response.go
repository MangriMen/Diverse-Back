package responses

// UploadDataResponseBody includes the id of uploaded data.
type UploadDataResponseBody struct {
	BaseResponseBody
	// required: true
	ID string `json:"id"`
}

// UploadDataResponse represent the response retrived on upload data request.
// swagger:response
type UploadDataResponse struct {
	// in: body
	Body UploadDataResponseBody
}
