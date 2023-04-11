// Package posthelpers provides functionality to convert post
// from DB to response variant.
package posthelpers

import (
	"fmt"
	"math"
	"time"

	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/MangriMen/Diverse-Back/internal/helpers/userhelpers"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/MangriMen/Diverse-Back/internal/parameters"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// PreparePostToSend prepares a post object for sending by fetching additional data from the database
// such as the user associated with the post and the comments associated with the post.
// The function then returns the prepared post object.
func PreparePostToSend(post models.DBPost, db *database.Queries) models.Post {
	preparedPost := post.ToPost()

	user, err := db.GetUser(post.UserID)
	if err == nil {
		preparedPost.User = helpers.Ptr(user.ToUser())
	}

	comments, err := db.GetComments(
		post.ID,
		&parameters.CommentsFetchRequestQuery{
			Count:                    configs.PostFetchCommentCount,
			LastSeenCommentCreatedAt: time.Now(),
		},
	)
	if err == nil {
		preparedPost.Comments = lo.Map(
			comments,
			func(item models.DBComment, index int) models.Comment {
				return PrepareCommentToPost(item, db)
			},
		)
	}

	return preparedPost
}

// PrepareCommentToPost prepares a comment object for inclusion in a post object
// by fetching additional data from the database
// such as the user associated with the comment.
// The function then returns the prepared comment object.
func PrepareCommentToPost(comment models.DBComment, db *database.Queries) models.Comment {
	preparedComment := comment.ToComment()

	user, err := db.GetUser(comment.UserID)
	if err == nil {
		preparedComment.User = helpers.Ptr(user.ToUser())
	}

	return preparedComment
}

// GenerateFilter TODO
func GenerateFilter(
	userID uuid.UUID,
	postsFetchRequestQuery *parameters.PostsFetchRequestQuery,
	db *database.Queries,
) (string, error) {
	var postFromCondition string

	rawRelationStatus, err := db.GetRelationStatus(&parameters.RelationGetStatusParams{
		UserIDParams: parameters.UserIDParams{User: userID},
		RelationUserIDParams: parameters.RelationUserIDParams{
			RelationUser: postsFetchRequestQuery.UserID,
		},
	})
	if err != nil {
		return "", err
	}

	relationStatus := userhelpers.PrepareRelationStatusToSend(rawRelationStatus)

	switch postsFetchRequestQuery.Type {
	case parameters.Subscriptions:
		if relationStatus["blocked"] {
			return "", fmt.Errorf("can't get posts, blocked by user")
		}

		rawRelations, rawRelationsErr := db.GetRelations(
			userID,
			&parameters.RelationGetRequestQuery{
				Type:  models.Following,
				Count: math.MaxInt64,
			},
		)
		if rawRelationsErr != nil {
			return "", rawRelationsErr
		}

		relations := lo.Map(
			rawRelations,
			func(item models.DBRelation, index int) models.Relation {
				return *userhelpers.PrepareRelationToSend(item, db)
			},
		)

		for _, relation := range relations {
			postFromCondition += fmt.Sprintf("AND user_id='%s'", relation.RelationUser.ID.String())
		}
	case parameters.User:
		if relationStatus["blocked"] {
			return "", fmt.Errorf("can't get posts, blocked by user")
		}

		postFromCondition = fmt.Sprintf("AND user_id='%s'", postsFetchRequestQuery.UserID)
	case parameters.All:
		postFromCondition = ""
	}

	return postFromCondition, nil
}
