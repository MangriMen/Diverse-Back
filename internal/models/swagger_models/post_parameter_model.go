package swagger_models

import "github.com/google/uuid"

// swagger:parameters getPost updatePost deletePost addComment updateComment deleteComment
type PostIdParameter struct {
	// in: path
	// required: true
	Post uuid.UUID `json:"post"`
}

// swagger:parameters updateComment deleteComment
type PostCommentIdParameter struct {
	// in: path
	// required: true
	Comment uuid.UUID `json:"comment"`
}
