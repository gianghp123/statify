package project

import (
	"github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/modules/project/service"
	"github.com/gin-gonic/gin"
)

type ProjectModule struct {
	controller *ProjectController
}

func NewProjectModule(repo repository.IProjectRepository) *ProjectModule {
	svc := service.NewProjectService(repo)
	ctrl := NewProjectController(svc)

	return &ProjectModule{
		controller: ctrl,
	}
}

func (m *ProjectModule) RegisterRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		projects.GET("", m.controller.ListProjects)
		projects.POST("", m.controller.CreateProject)
		projects.GET("/:project_id", m.controller.GetProject)
	}
}
