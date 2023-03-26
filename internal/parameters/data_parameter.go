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
