package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gianghp/statify/internal/core"
	eventmessage "github.com/gianghp/statify/internal/core/event-message"
	"github.com/gianghp/statify/internal/core/sse"
	"github.com/gianghp/statify/internal/modules/analytics/service"
	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	broker  *sse.Broker
	service *service.AnalyticsService
}

func NewAnalyticsController(broker *sse.Broker, service *service.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{
		broker:  broker,
		service: service,
	}
}

func (c *AnalyticsController) StreamAnalyticsEvents(ctx *gin.Context) {
	startTimeQuery := ctx.Query("start_time")
	endTimeQuery := ctx.Query("end_time")

	var startTime time.Time
	var endTime time.Time
	var err error

	if startTimeQuery != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeQuery)
		if err != nil {
			core.HandleApiError(ctx, core.BadRequestError("Invalid start time format"))
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -7)
	}

	if endTimeQuery != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeQuery)
		if err != nil {
			core.HandleApiError(ctx, core.BadRequestError("Invalid end time format"))
			return
		}
	} else {
		endTime = time.Now()
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		core.HandleApiError(ctx, core.InternalError("Failed to cast writer to flusher"))
		return
	}

	clientChan := make(chan string, 10)
	c.broker.Subscribe(sse.AnalyticsEvent, clientChan)

	go func() {
		<-ctx.Done()
		c.broker.Unsubscribe(sse.AnalyticsEvent, clientChan)
	}()

	// Stream events
	for msg := range clientChan {
		var jsonMsg eventmessage.AnalyticsMessage

		err := json.Unmarshal([]byte(msg), &jsonMsg)
		if err != nil {
			core.HandleApiError(ctx, core.InternalError(err.Error()))
			return
		}

		metrics, err := c.service.GetComprehensiveMetrics(ctx, jsonMsg.ProjectID, startTime, endTime)
		if err != nil {
			core.HandleApiError(ctx, err)
			return
		}

		payload, err := json.Marshal(metrics)
		if err != nil {
			core.HandleApiError(ctx, core.InternalError(err.Error()))
			return
		}

		_, err = fmt.Fprintf(ctx.Writer, "data: %s\n\n", payload)
		if err != nil {
			core.HandleApiError(ctx, core.InternalError(err.Error()))
			return
		}
		flusher.Flush()
	}
}

func (c *AnalyticsController) GetComprehensiveMetrics(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")
	projectId, err := strconv.ParseUint(projectIDStr, 10, 32)

	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID format"))
		return
	}

	startTimeQuery := ctx.Query("start_time")
	endTimeQuery := ctx.Query("end_time")

	var startTime time.Time
	var endTime time.Time

	if startTimeQuery != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeQuery)
		if err != nil {
			core.HandleApiError(ctx, core.BadRequestError("Invalid start time format"))
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -7)
	}

	if endTimeQuery != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeQuery)
		if err != nil {
			core.HandleApiError(ctx, core.BadRequestError("Invalid end time format"))
			return
		}
	} else {
		endTime = time.Now()
	}

	metrics, err := c.service.GetComprehensiveMetrics(ctx, uint(projectId), startTime, endTime)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Metrics retrieved successfully", metrics))
}
