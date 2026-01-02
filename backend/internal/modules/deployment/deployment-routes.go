package deployment

import "github.com/gin-gonic/gin"

type DeploymentRoutes struct {
	controller *DeploymentController
}

func NewDeploymentRoutes(controller *DeploymentController) *DeploymentRoutes {
	return &DeploymentRoutes{controller: controller}
}

func (r *DeploymentRoutes) SetupDeploymentRoutes(rg *gin.RouterGroup) {
	deployments := rg.Group("/projects/:project_id/deployments")
	{
		deployments.GET("", r.controller.GetHistory)
		deployments.POST("", r.controller.CreateDeployment)
		deployments.GET("/:id", r.controller.GetStatus)
	}
}
