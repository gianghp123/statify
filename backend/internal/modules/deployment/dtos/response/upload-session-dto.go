package response

import "time"

type UploadSessionDto struct {
	ID           uint      `json:"id"`
	ProjectID    uint      `json:"project_id"`
	UploadKey    string    `json:"upload_key"`
	OutputPrefix string    `json:"output_prefix"`
	PresignedUrl string    `json:"presigned_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}
