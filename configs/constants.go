// Package configs provides base constants
package configs

import "time"

// PasswordEncryptCost is computational cost for bcrypt algorithm.
const PasswordEncryptCost = 12

// PostEditTimeSinceCreated is a threshold value for determining whether a post can be edited or not.
const PostEditTimeSinceCreated = 24 * time.Hour

// PostCommentEditTimeSinceCreated is a threshold value for determining whether a comment can be edited or not.
const PostCommentEditTimeSinceCreated = 24 * time.Hour

// PostFetchCommentCount specifies the maximum number of comments to first time fetch a post.
const PostFetchCommentCount = 20

// DataPath specifies the root path for image files, video files, etc.
const DataPath = "/var/lib/backend-data/"

// BodyLimit is limit for body size in bits.
const BodyLimit = 1024 * 1024 * 1024 * 8

// Constants for MIME base types.
const (
	MIMEBaseImage = "image"
)

// GetAllowedMIMEBaseTypes returns types allowed to upload on server.
func GetAllowedMIMEBaseTypes() []string {
	return []string{MIMEBaseImage}
}

// WebpQuality is webp quality for saved images.
const WebpQuality = 85

// Constants for postgres errors.
const (
	DBDuplicateError = "23505"
)

// Constants for response.
const (
	UserNotFoundError      = "User with this ID not found"
	UsersNotFoundError     = "Users not found"
	UserAlreadyExistsError = "User with this email or username already exists"
	/* #nosec */
	WrongEmailOrPasswordError = "Wrong email or password"

	PostNotFoundError  = "Post with this ID not found"
	PostsNotFoundError = "Posts not found"

	CommentNotFoundError  = "Comment with this ID not found"
	CommentsNotFoundError = "Comments not found"

	ForbiddenError = "Not enough permission"

	CantEditAfterErrorFormat = "Can't edit %s after %s"
)

// TestResponseTimeout is response timeout for tests.
const TestResponseTimeout = 10000
