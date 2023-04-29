package parameters

// RequestParams is interface to union all request parameters in one type.
type RequestParams interface {
	UserIDParams |
		UsernameIDParams |
		RelationGetStatusParams |
		PostIDParams |
		PostCommentIDParams |
		CommentAddRequestParams |
		GetDataRequestParams
}

// RequestQuery is interface to union all request queries in one type.
type RequestQuery interface {
	PostsFetchCountRequestQuery |
		PostsFetchRequestQuery |
		CommentsFetchRequestQuery |
		RelationGetCountRequestQuery |
		RelationGetRequestQuery |
		RelationAddDeleteRequestQuery |
		GetDataRequestQuery
}

// RequestBody is interface to union all request body in one type.
type RequestBody interface {
	LoginRequestBody |
		RegisterRequestBody |
		UserUpdateRequestBody |
		PostCreateRequestBody |
		PostUpdateRequestBody |
		CommentAddRequestBody |
		CommentUpdateRequestBody
}
