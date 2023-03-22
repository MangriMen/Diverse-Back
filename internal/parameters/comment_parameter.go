package parameters

import "github.com/google/uuid"

type PostCommentIdParams struct {
	PostIdParams

	// in: path
	// required: true
	Comment uuid.UUID `params:"comment" json:"comment" validate:"required"`
}

// swagger:parameters deleteComment
type PostCommentIdRequest struct {
	PostCommentIdParams
}

type CommentAddRequestParams struct {
	PostIdParams
}

type CommentAddRequestBody struct {
	// required: true
	Content string `json:"content" validate:"required"`
}

// swagger:parameters addComment
type CommentAddRequest struct {
	CommentAddRequestParams

	// in: body
	// required: true
	Body CommentAddRequestBody
}

type CommentUpdateRequestBody struct {
	// required: true
	Content string `json:"content" validate:"required"`
}

// swagger:parameters updateComment
type CommentUpdateRequest struct {
	PostCommentIdParams

	// in: body
	// required: true
	Body CommentUpdateRequestBody
}
