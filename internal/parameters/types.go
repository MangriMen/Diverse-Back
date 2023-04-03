package parameters

// RequestParams is interface to union all request parameters in one type.
type RequestParams interface {
	UserIDParams |
		RelationAddRequestParams |
		PostIDParams |
		PostCommentIDParams |
		CommentAddRequestParams
}

// RequestQuery is interface to union all request queries in one type.
type RequestQuery interface {
	PostsFetchRequestQuery |
		CommentsFetchRequestQuery |
		RelationGetRequestQuery
}

// RequestBody is interface to union all request body in one type.
type RequestBody interface {
	LoginRequestBody |
		RegisterRequestBody |
		UserUpdateRequestBody |
		RelationAddRequestBody |
		PostCreateRequestBody |
		PostUpdateRequestBody |
		CommentAddRequestBody |
		CommentUpdateRequestBody
}
