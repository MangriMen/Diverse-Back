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
	UserNotFoundError      = "user with this ID not found"
	UsersNotFoundError     = "users not found"
	UserAlreadyExistsError = "user with this email or username already exists"
	/* #nosec */
	WrongEmailOrPasswordError = "wrong email or password"
	WrongPassword             = "wrong password"
	UserBlocked               = "blocked by user"
	RelationsGetError         = "relations getting error"

	PostNotFoundError  = "post with this ID not found"
	PostsNotFoundError = "posts not found"
	PostsInvalidFilter = "invalid filter option"

	CommentNotFoundError  = "comment with this ID not found"
	CommentsNotFoundError = "comments not found"

	ForbiddenError = "not enough permission"

	CantEditAfterErrorFormat = "can't edit %s after %s"
)

// TestResponseTimeout is response timeout for tests.
const TestResponseTimeout = 10000

// MockDSN is custom dsn for mock database.
const MockDSN = "mock_dsn"

// MockLogFilename is filename for sql queries to mock db.
const MockLogFilename = "mock_db.log"

// TempDirectory is dir for unnecessary runtime logs and etc.
const TempDirectory = "tmp"
