package response

type DeploymentStatusDTO struct {
	ID     uint    `json:"id"`
	Status string  `json:"status"`
	Reason *string `json:"reason,omitempty"`
}
