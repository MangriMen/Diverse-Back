package responses

import "bytes"

// UploadDataResponseBody includes the relative path to uploaded data.
type UploadDataResponseBody struct {
	BaseResponseBody

	// required: true
	Path string `json:"path" validate:"required"`
}

// UploadDataResponse contains the uploaded data info.
// swagger:response
type UploadDataResponse struct {
	// in: body
	Body UploadDataResponseBody
}

// GetDataResponse contains the binary data of requested type.
// swagger:response
type GetDataResponse struct {
	// in: body
	// swagger:file
	Body *bytes.Buffer
}
