package core

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

func BadRequestError() *ApiError {
	return NewApiError(400, "Bad Request")
}

func NotFoundError() *ApiError {
	return NewApiError(404, "Not Found")
}

func UnauthorizedError() *ApiError {
	return NewApiError(401, "Unauthorized")
}

func ForbiddenError() *ApiError {
	return NewApiError(403, "Forbidden")
}

func InternalError() *ApiError {
	return NewApiError(500, "Internal Server Error")
}

func NotImplementedError() *ApiError {
	return NewApiError(501, "Not Implemented")
}

func BadGatewayError() *ApiError {
	return NewApiError(502, "Bad Gateway")
}

func ServiceUnavailableError() *ApiError {
	return NewApiError(503, "Service Unavailable")
}
