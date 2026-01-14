package response

import (
	"time"

	"github.com/gianghp/statify/internal/core/enums"
)

type ProjectDto struct {
	ID        uint                   `json:"id"`
	Name      string                 `json:"name"`
	Subdomain string                 `json:"subdomain"`
	CreatedAt time.Time              `json:"created_at"`
	URL       string                 `json:"url"`
	Status    enums.DeploymentStatus `json:"status"`
	UpdatedAt time.Time              `json:"updated_at"`
}
