package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func BadRequestError(message ...string) *ApiError {
	msg := "Bad Request"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return NewApiError(http.StatusBadRequest, msg)
}

func NotFoundError() *ApiError {
	return NewApiError(http.StatusNotFound, "Not Found")
}

func ConflictError(messages ...string) *ApiError {
	msg := "Conflict"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	return NewApiError(http.StatusConflict, msg)
}

func UnauthorizedError() *ApiError {
	return NewApiError(http.StatusUnauthorized, "Unauthorized")
}

func ForbiddenError() *ApiError {
	return NewApiError(http.StatusForbidden, "Forbidden")
}

func InternalError() *ApiError {
	return NewApiError(http.StatusInternalServerError, "Internal Server Error")
}

func NotImplementedError() *ApiError {
	return NewApiError(http.StatusNotImplemented, "Not Implemented")
}

func BadGatewayError() *ApiError {
	return NewApiError(http.StatusBadGateway, "Bad Gateway")
}

func ServiceUnavailableError() *ApiError {
	return NewApiError(http.StatusServiceUnavailable, "Service Unavailable")
}

func HandleApiError(ctx *gin.Context, err error) {
	if apiErr, ok := err.(*ApiError); ok {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, InternalError())
}
