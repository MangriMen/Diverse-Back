package parameters

// RequestParams is interface to union all request parameters in one type.
type RequestParams interface {
	UserIDParams |
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
		GetDataRequestQuery
}

// RequestBody is interface to union all request body in one type.
type RequestBody interface {
	LoginRequestBody |
		RegisterRequestBody |
		UserUpdateRequestBody |
		RelationAddDeleteRequestBody |
		PostCreateRequestBody |
		PostUpdateRequestBody |
		CommentAddRequestBody |
		CommentUpdateRequestBody
}
