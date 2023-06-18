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

// DataTypeParams includes the type of the data.
type DataTypeParams struct {
	// in: path
	// required: true
	Type string `params:"type" json:"type" validate:"required"`
}

// DataIDParams includes the id of the image.
type DataIDParams struct {
	// in: path
	// required: true
	ID string `params:"id" json:"id" validate:"required"`
}

// DataImageIDRequest is used to represent a request thet requires a
// image id to get image.
// swagger:parameters getImageRaw
type DataImageIDRequest struct {
	DataIDParams
}

// GetDataRequestParams TODO.
type GetDataRequestParams struct {
	DataTypeParams
	DataIDParams
}

// GetDataRequestQuery TODO.
type GetDataRequestQuery struct {
	// in: query
	// required: true
	Width *int `query:"width" json:"width"`

	// in: query
	// required: true
	Height *int `query:"height" json:"height"`
}

// GetDataRequest TODO.
// swagger:parameters getData
type GetDataRequest struct {
	GetDataRequestParams
	GetDataRequestQuery
}
