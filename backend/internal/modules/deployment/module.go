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
	rg.GET("/serve-files/*file_name", m.controller.ServeFiles)
	deployments := rg.Group("/projects/:project_id/deployments")
	{
		deployments.GET("", authMiddleware, m.controller.GetHistory)
		deployments.POST("", authMiddleware, m.controller.CreateDeployment)
		deployments.GET("/:id", authMiddleware, m.controller.GetStatus)
		deployments.DELETE("/:id", authMiddleware, m.controller.DeleteDeployment)
		deployments.PUT("/:id/live", authMiddleware, m.controller.TurnDeploymentLive)
		deployments.PUT("/:id/offline", authMiddleware, m.controller.TurnDeploymentOffline)
	}
	globalDeployments := rg.Group("/deployments")
	{
		globalDeployments.GET("", authMiddleware, m.controller.GetGlobalDeploymentHistory)
		globalDeployments.PUT("/:id/toggle-spa-mode", authMiddleware, m.controller.ToggleIsSPAMode)
		globalDeployments.DELETE("/:id", authMiddleware, m.controller.DeleteDeployment)
	}
}
