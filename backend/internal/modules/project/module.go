package project

import (
	"github.com/gin-gonic/gin"
)

type ProjectModule struct {
	controller *ProjectController
}

func NewProjectModule(controller *ProjectController) *ProjectModule {
	return &ProjectModule{
		controller: controller,
	}
}

func (m *ProjectModule) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	projects := rg.Group("/projects")
	{
		projects.GET("", authMiddleware, m.controller.ListProjects)
		projects.POST("", authMiddleware, m.controller.CreateProject)
		projects.GET("/:project_id", authMiddleware, m.controller.GetProject)
		projects.PUT("/:project_id", authMiddleware, m.controller.UpdateProject)
		projects.DELETE("/:project_id", authMiddleware, m.controller.DeleteProject)
	}
}
