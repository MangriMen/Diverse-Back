package helpers

import (
	"github.com/MangriMen/Diverse-Back/api/database"
	"github.com/MangriMen/Diverse-Back/internal/models"
	"github.com/samber/lo"
)

func PreparePostToSend(post models.DBPost, db *database.Queries) models.Post {
	preparedPost := post.ToPost()

	user, err := db.GetUser(post.UserId)
	if err == nil {
		preparedPost.User = &user.BaseUser
	}

	comments, err := db.GetComments(post.Id)
	if err == nil {
		preparedPost.Comments = lo.Map(comments, func(item models.DBComment, index int) models.Comment {
			return PrepareCommentToPost(item, db)
		})
	}

	return preparedPost
}

func PrepareCommentToPost(comment models.DBComment, db *database.Queries) models.Comment {
	preparedComment := comment.ToComment()

	user, err := db.GetUser(comment.UserId)
	if err == nil {
		preparedComment.User = &user.BaseUser
	}

	return preparedComment
}
