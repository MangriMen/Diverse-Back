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
func PreparePostToSend(post models.DBPost, userID uuid.UUID, db *database.Queries) models.Post {
	preparedPost := post.ToPost()

	isLikedByRequester, err := db.GetPostIsLiked(post.ID, userID)
	if err == nil {
		preparedPost.LikedByMe = isLikedByRequester
	}

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
				return PrepareCommentToPost(item, userID, db)
			},
		)
	}

	return preparedPost
}

// PrepareCommentToPost prepares a comment object for inclusion in a post object
// by fetching additional data from the database such as the user associated with the comment.
func PrepareCommentToPost(
	comment models.DBComment,
	userID uuid.UUID,
	db *database.Queries,
) models.Comment {
	preparedComment := comment.ToComment()

	isLikedByRequester, err := db.GetCommentIsLiked(comment.ID, userID)
	if err == nil {
		preparedComment.LikedByMe = isLikedByRequester
	}

	user, err := db.GetUser(comment.UserID)
	if err == nil {
		preparedComment.User = helpers.Ptr(user.ToUser())
	}

	return preparedComment
}

// GenerateFilter generates a filter for SQL query to fetch posts by the specified parameters.
func GenerateFilter(
	userID uuid.UUID,
	postFetchType parameters.PostFetchType,
	postsAuthorID uuid.UUID,
	db *database.Queries,
) (string, error) {
	const conditionFormatString = "AND user_id='%s'"

	switch postFetchType {
	case parameters.Subscriptions:
		rawRelations, rawRelationsErr := db.GetRelations(
			userID,
			&parameters.RelationGetRequestQuery{
				Type:                      models.Following,
				Count:                     math.MaxInt64,
				LastSeenRelationCreatedAt: time.Now(),
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

		condition := ""
		for _, relation := range relations {
			condition += fmt.Sprintf(conditionFormatString, relation.RelationUser.ID.String())
		}
		return condition, nil
	case parameters.User:
		rawRelationStatus, err := db.GetRelationStatus(&parameters.RelationGetStatusParams{
			UserIDParams: parameters.UserIDParams{User: userID},
			RelationUserIDParams: parameters.RelationUserIDParams{
				RelationUser: postsAuthorID,
			},
		})
		if err != nil {
			return "", err
		}

		relationStatus := userhelpers.PrepareRelationStatusToSend(userID, rawRelationStatus)

		if relationStatus["blocked"] {
			return "", fmt.Errorf(configs.UserBlocked)
		}

		return fmt.Sprintf(conditionFormatString, postsAuthorID), nil
	case parameters.All:
		return "", nil
	default:
		return "", fmt.Errorf(configs.PostsInvalidFilter)
	}
}
