package response

import "time"

type DeploymentDto struct {
	ID                 uint       `json:"id"`
	ProjectID          uint       `json:"project_id"`
	Status             string     `json:"status"`
	OutputPrefix       *string    `json:"output_prefix"`
	SourceZipObjectKey *string    `json:"source_zip_object_key"`
	ValidationError    *string    `json:"validation_error"`
	CreatedAt          time.Time  `json:"created_at"`
	FinishedAt         *time.Time `json:"finished_at"`
}
