package deployment

import (
	"github.com/gianghp/statify/internal/modules/analytics/wrapper"
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

func (m *DeploymentModule) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc, analyzerWrapper *wrapper.AnalyzerWrapper) {
	rg.GET("/serve-files/*file_name", analyzerWrapper.Wrap(m.controller.ServeFiles))
	deployments := rg.Group("/projects/:project_id/deployments")
	{
		deployments.GET("", authMiddleware, m.controller.GetHistory)
		deployments.GET("/upload-session", authMiddleware, m.controller.CreateUploadSession)
		deployments.PUT("/confirm/:upload_session_id", authMiddleware, m.controller.ConfirmCreateDeployment)
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
		globalDeployments.GET("/stream-status", authMiddleware, m.controller.StreamDeploymentStatus)
	}
}
