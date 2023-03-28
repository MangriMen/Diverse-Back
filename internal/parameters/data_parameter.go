package parameters

// UploadDataRequestForm includes the file being loaded.
type UploadDataRequestForm struct {
	// in: formData
	// swagger:file
	File interface{} `json:"file"`
}

// UploadDataRequest is used for upload a new file.
// swagger:parameters uploadData
type UploadDataRequest struct {
	UploadDataRequestForm
}

// DataImageIDParams includes the id of the image.
type DataImageIDParams struct {
	// in: path
	// required: true
	Image string `params:"image" json:"image" validate:"required"`
}

// DataImageIDRequest is used to represent a request thet requires a
// image id to get image.
// swagger:parameters getImageRaw
type DataImageIDRequest struct {
	DataImageIDParams
}
