package response

import (
	"time"

	projecResponse "github.com/gianghp/statify/internal/modules/project/dtos/response"
)

type DeploymentDto struct {
	ID              uint                       `json:"id"`
	Project         *projecResponse.ProjectDto `json:"project,omitempty"`
	Status          string                     `json:"status"`
	ValidationError *string                    `json:"validation_error,omitempty"`
	CreatedAt       time.Time                  `json:"created_at"`
	FinishedAt      *time.Time                 `json:"finished_at,omitempty"`
}
