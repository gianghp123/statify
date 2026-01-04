package deployment

import (
	"github.com/gin-gonic/gin"
)

type DeploymentModule struct {
	controller *DeploymentController
}

func NewDeploymentModule(controller *DeploymentController) *DeploymentModule {
	return &DeploymentModule{
		controller: controller,
	}
}

func (m *DeploymentModule) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	deployments := rg.Group("/projects/:project_id/deployments")
	{
		deployments.GET("", authMiddleware, m.controller.GetHistory)
		deployments.POST("", authMiddleware, m.controller.CreateDeployment)
		deployments.GET("/:id", authMiddleware, m.controller.GetStatus)
	}
}
