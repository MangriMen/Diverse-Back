package responses

// UploadDataResponseBody includes the id of uploaded data.
type UploadDataResponseBody struct {
	BaseResponseBody
	// required: true
	Path string `json:"path"`
}

// UploadDataResponse represent the response retrived on upload data request.
// swagger:response
type UploadDataResponse struct {
	// in: body
	Body UploadDataResponseBody
}
