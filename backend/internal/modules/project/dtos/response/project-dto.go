package response

import (
	"time"

	"github.com/gianghp/statify/internal/core/enums"
)

type ProjectDto struct {
	ID                  uint                   `json:"id"`
	Name                string                 `json:"name"`
	Subdomain           string                 `json:"subdomain"`
	UserID              uint                   `json:"user_id"`
	CurrentDeploymentID *uint                  `json:"current_deployment_id"`
	CreatedAt           time.Time              `json:"created_at"`
	URL                 string                 `json:"url"`
	Status              enums.DeploymentStatus `json:"status"`
	// LastCommit          string    `json:"last_commit"`
	// LastCommitHash      string    `json:"last_commit_hash"`
	UpdatedAt time.Time `json:"updated_at"`
}
