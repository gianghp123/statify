package project

import "github.com/gin-gonic/gin"

type ProjectRoutes struct {
	controller *ProjectController
}

func NewProjectRoutes(controller *ProjectController) *ProjectRoutes {
	return &ProjectRoutes{controller: controller}
}

func (r *ProjectRoutes) SetupProjectRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		projects.GET("", r.controller.ListProjects)
		projects.POST("", r.controller.CreateProject)
		projects.GET("/:id", r.controller.GetProject)
	}
}
