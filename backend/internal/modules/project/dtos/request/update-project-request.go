package request

type UpdateProjectRequest struct {
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
}
