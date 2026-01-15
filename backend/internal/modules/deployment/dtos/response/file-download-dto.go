package response

import "io"

type FileDownloadDto struct {
	// The actual data stream (MinIO Object)
	// We use ReadCloser so the Controller can close the connection after serving
	Stream io.ReadCloser

	// Metadata required by c.DataFromReader
	Size        int64
	ContentType string

	// Headers for the response (ETag, Cache-Control, etc.)
	Headers map[string]string

	// If true, Controller should send 304 and NOT read the Stream
	NotModified bool

	StatusCode int
}
