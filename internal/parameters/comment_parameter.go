// Package parameters provides base parameter models for swagger and work with REST API
package parameters

import (
	"time"

	"github.com/google/uuid"
)

// PostCommentIDParams includes the id of the post and id of the comment.
type PostCommentIDParams struct {
	PostIDParams

	// in: path
	// required: true
	Comment uuid.UUID `params:"comment" json:"comment" validate:"required"`
}

// PostCommentIDRequest is used to represent a request thet requires a
// post id and comment id parameters, such as deleting comment.
// swagger:parameters deleteComment likeComment unlikeComment
type PostCommentIDRequest struct {
	PostCommentIDParams
}

// CommentAddRequestParams include the ID of the post where you want to create a comment.
type CommentAddRequestParams struct {
	PostIDParams
}

// CommentAddRequestBody includes the content of the comment.
type CommentAddRequestBody struct {
	// required: true
	Content string `json:"content" validate:"required"`
}

// CommentAddRequest is used for adding a new comment to post.
// It includes the new comment content.
// swagger:parameters addComment
type CommentAddRequest struct {
	CommentAddRequestParams

	// in: body
	// required: true
	Body CommentAddRequestBody
}

// CommentUpdateRequestBody includes the content of the comment.
type CommentUpdateRequestBody struct {
	// required: true
	Content string `json:"content" validate:"required"`
}

// CommentUpdateRequest is used for updating a comment of the post.
// It includes the new comment content.
// swagger:parameters updateComment
type CommentUpdateRequest struct {
	PostCommentIDParams

	// in: body
	// required: true
	Body CommentUpdateRequestBody
}

// CommentsFetchRequestQuery includes the ID and creation time of the last seen comment,
// as well as a count of the number of comments to retrieve.
type CommentsFetchRequestQuery struct {
	// in: query
	LastSeenCommentID uuid.UUID `query:"last_seen_comment_id" json:"last_seen_comment_id" validate:"uuid"`

	//nolint:lll
	// in: query
	LastSeenCommentCreatedAt time.Time `query:"last_seen_comment_created_at" json:"last_seen_comment_created_at" validate:"uuid,required_with=last_seen_comment_id"`

	// in: query
	// required: true
	// min: 1
	// max: 50
	Count int `query:"count" json:"count" validate:"required,min=1,max=50"`
}

// CommentsFetchRequest is a struct that encapsulates a query used to fetch comments.
// swagger:parameters getComments
type CommentsFetchRequest struct {
	PostIDParams

	CommentsFetchRequestQuery
}
