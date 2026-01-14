package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationConfig extracts page and limit from query parameters.
// Defaults: page = 1, limit = 10.
func GetPaginationConfig(ctx *gin.Context) (page int, limit int) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Optional: Cap the limit to prevent over-fetching
	if limit > 100 {
		limit = 100
	}

	return page, limit
}
