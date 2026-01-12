package core

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
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

func UnauthorizedError(messages ...string) *ApiError {
	msg := "Unauthorized"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	return NewApiError(http.StatusUnauthorized, msg)
}

func ForbiddenError(messages ...string) *ApiError {
	msg := "Forbidden"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	return NewApiError(http.StatusForbidden, msg)
}

func InternalError(message ...string) *ApiError {
	msg := "Internal Server Error"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return NewApiError(http.StatusInternalServerError, msg)
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

func ParseDatabaseError(err error) *ApiError {
	if err == nil {
		return nil
	}
	if IsRecordNotFoundError(err) {
		return NotFoundError()
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ConflictError("Record already exists")
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return BadRequestError("Foreign key violation")
	}
	return InternalError(err.Error())
}

func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func ParseMinioError(err error) *ApiError {
	if err == nil {
		return nil
	}

	// Convert the generic error to a Minio-specific ErrorResponse
	minioErr := minio.ToErrorResponse(err)

	switch minioErr.Code {
	case "NoSuchKey", "NoSuchBucket":
		return NotFoundError()
	case "AccessDenied":
		return ForbiddenError("Access to storage was denied")
	default:
		// Handle network/connection issues or generic internal errors
		return InternalError(err.Error())
	}
}

func HandleApiError(ctx *gin.Context, err error) {
	ctx.Error(err)

	// If headers were already written, don't try to send another response
	if ctx.Writer.Written() {
		return
	}

	if apiErr, ok := err.(*ApiError); ok {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, InternalError())
}
