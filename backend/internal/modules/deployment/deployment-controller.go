package deployment

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/sse"
	"github.com/gianghp/statify/internal/modules/analytics/wrapper"
	"github.com/gianghp/statify/internal/modules/deployment/service"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gin-gonic/gin"
)

type DeploymentController struct {
	service *service.DeploymentService
	broker  *sse.Broker
}

func NewDeploymentController(service *service.DeploymentService, broker *sse.Broker) *DeploymentController {
	return &DeploymentController{service: service, broker: broker}
}

func (c *DeploymentController) GetGlobalDeploymentHistory(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	page, limit := utils.GetPaginationConfig(ctx)
	deployments, err := c.service.GetGlobalDeploymentHistory(ctx, userID, page, limit)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Global deployment history retrieved successfully", deployments.Entities, &deployments.Pagination))
}

func (c *DeploymentController) GetHistory(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	page, limit := utils.GetPaginationConfig(ctx)
	deployments, err := c.service.GetHistory(ctx, userID, uint(projectID), page, limit)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Deployment history retrieved successfully", deployments.Entities, &deployments.Pagination))
}

func (c *DeploymentController) CreateUploadSession(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	uploadSessionDto, err := c.service.CreatePresignedUrl(ctx, userID, uint(projectID))

	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "Upload session created successfully", uploadSessionDto))
}

func (c *DeploymentController) GetStatus(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	deploymentDto, err := c.service.GetDeploymentByID(ctx, userID, uint(id))
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Deployment status retrieved successfully", deploymentDto))
}

func (c *DeploymentController) ServeFiles(ctx *gin.Context) (*wrapper.StaticServeResult, error) {
	host := strings.Split(ctx.Request.Host, ".")[0]

	clientEtag := ctx.GetHeader("If-None-Match")
	fileName := ctx.Param("file_name")

	if fileName == "" || fileName == "/" {
		fileName = "index.html"
	}

	// Optional: strip the leading slash if it exists
	fileName = strings.TrimPrefix(fileName, "/")

	log.Println("host", host, "fileName", fileName, "clientEtag", clientEtag)
	fileDTO, err := c.service.GetCurrentDeploymentFilesByProjectSubdomain(ctx, host, fileName, clientEtag)

	if err != nil {
		core.HandleApiError(ctx, err)
		return nil, err
	}

	if fileDTO.Stream != nil {
		defer fileDTO.Stream.Close()
	}

	if fileDTO.NotModified {
		ctx.Status(http.StatusNotModified)
		return &wrapper.StaticServeResult{
			DeploymentID: fileDTO.DeploymentID,
			ProjectID:    fileDTO.ProjectID,
			StatusCode:   http.StatusNotModified,
		}, nil
	}

	ctx.DataFromReader(fileDTO.StatusCode, fileDTO.Size, fileDTO.ContentType, fileDTO.Stream, fileDTO.Headers)
	return &wrapper.StaticServeResult{
		DeploymentID: fileDTO.DeploymentID,
		ProjectID:    fileDTO.ProjectID,
		StatusCode:   fileDTO.StatusCode,
		BytesServed:  fileDTO.Size,
	}, nil
}

func (c *DeploymentController) DeleteDeployment(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	err = c.service.DeleteDeployment(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse[any](http.StatusOK, "Deployment deleted successfully", nil))
}

func (c *DeploymentController) TurnDeploymentLive(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	err = c.service.TurnDeploymentLive(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, core.NewApiResponse[any](http.StatusNoContent, "Deployment turned live successfully", nil))
}

func (c *DeploymentController) TurnDeploymentOffline(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	err = c.service.TurnDeploymentOffline(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, core.NewApiResponse[any](http.StatusNoContent, "Deployment turned offline successfully", nil))
}

func (c *DeploymentController) ToggleIsSPAMode(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	err = c.service.ToggleIsSPAMode(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, core.NewApiResponse[any](http.StatusNoContent, "Deployment is SPA mode toggled successfully", nil))
}

func (c *DeploymentController) StreamDeploymentStatus(ctx *gin.Context) {
	//Set SSE headers
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	clientChan := make(chan string, 10)
	c.broker.Subscribe(sse.DeploymentStatusEvent, clientChan)

	go func() {
		<-ctx.Done()
		c.broker.Unsubscribe(sse.DeploymentStatusEvent, clientChan)
	}()

	// Stream events
	for msg := range clientChan {
		_, err := fmt.Fprintf(ctx.Writer, "data: %s\n\n", msg)
		if err != nil {
			return
		}
		flusher.Flush()
	}
}
