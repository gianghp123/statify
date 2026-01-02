package deployment

import (
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/deployment/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeploymentModule struct {
	controller *DeploymentController
}

func NewDeploymentModule(db *gorm.DB) *DeploymentModule {
	repo := repository.NewDeploymentRepository(db)
	svc := service.NewDeploymentService(repo)
	ctrl := NewDeploymentController(svc)

	return &DeploymentModule{
		controller: ctrl,
	}
}

func (m *DeploymentModule) RegisterRoutes(rg *gin.RouterGroup) {
	deployments := rg.Group("/projects/:project_id/deployments")
	{
		deployments.GET("", m.controller.GetHistory)
		deployments.POST("", m.controller.CreateDeployment)
		deployments.GET("/:id", m.controller.GetStatus)
	}
}
