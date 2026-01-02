package request

import "mime/multipart"

type CreateDeploymentRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
