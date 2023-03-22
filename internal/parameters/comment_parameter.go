package parameters

import (
	"time"

	"github.com/google/uuid"
)

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

type CommentsFetchRequestQuery struct {
	// in: query
	LastSeenCommentId uuid.UUID `query:"last_seen_comment_id" json:"last_seen_comment_id" validate:"uuid"`

	// in: query
	LastSeenCommentCreatedAt time.Time `query:"last_seen_comment_created_at" json:"last_seen_comment_created_at" validate:"uuid,required_with=last_seen_comment_id"`

	// in: query
	// required: true
	// min: 1
	// max: 50
	Count int `query:"count" json:"count" validate:"required,min=1,max=50"`
}

// swagger:parameters getPosts
type CommentsFetchRequest struct {
	CommentsFetchRequestQuery
}
