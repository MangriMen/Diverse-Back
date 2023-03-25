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
