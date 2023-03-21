package parameters

import "github.com/google/uuid"

// swagger:parameters deleteComment
type PostCommentIdParameter struct {
	PostIdParameter

	// in: path
	// required: true
	Comment uuid.UUID `params:"comment" validate:"required"`
}

type CommentAddParametersParams struct {
	PostIdParameter
}

type CommentAddParametersBody struct {
	// required: true
	Content string `json:"content" validate:"required"`
}

// swagger:parameters addComment
type CommentAddParameters struct {
	CommentAddParametersParams

	// in: body
	// required: true
	Body CommentAddParametersBody
}

// swagger:parameters updateComment
type CommentUpdateParameters struct {
	PostCommentIdParameter

	// in: body
	// required: true
	Body struct {
		// required: true
		Content string `json:"content" validate:"required"`
	} `json:"body"`
}
