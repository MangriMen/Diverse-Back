package responses

// ResponseBody is interface to union all response body in one type.
type ResponseBody interface {
	BaseResponseBody |
		GetUsersResponseBody |
		GetUserResponseBody |
		RegisterLoginUserResponseBody
}
