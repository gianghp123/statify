package response

import "time"

type ProjectDto struct {
	ID                  uint      `json:"id"`
	Name                string    `json:"name"`
	Subdomain           string    `json:"subdomain"`
	UserID              uint      `json:"user_id"`
	CurrentDeploymentID *uint     `json:"current_deployment_id"`
	CreatedAt           time.Time `json:"created_at"`
}
