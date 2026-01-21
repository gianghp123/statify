package wrapper

type StaticServeResult struct {
	DeploymentID uint
	ProjectID    uint
	StatusCode   int
	BytesServed  int64
}
