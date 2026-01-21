package analytics

import "github.com/gin-gonic/gin"

type AnalyticsModule struct {
	controller *AnalyticsController
}

func NewAnalyticsModule(controller *AnalyticsController) *AnalyticsModule {
	return &AnalyticsModule{
		controller: controller,
	}
}

func (m *AnalyticsModule) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	rg.GET("/projects/:project_id/stream-analytics", authMiddleware, m.controller.StreamAnalyticsEvents)
	rg.GET("/projects/:project_id/analytics", authMiddleware, m.controller.GetComprehensiveMetrics)
}
