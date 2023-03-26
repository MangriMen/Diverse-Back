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
const BodyLimit = 1024 * 1024 * 1024

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
