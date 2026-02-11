package core

import "time"

const (
	MaxDeploymentSize          = 50 * 1024 * 1024
	MaxFileCount               = 1000
	MaxProcessGoroutines       = 10
	MaxDeleteGoroutines        = 10
	MaxProjectDeleteGoroutines = 5
	MaxRetry                   = 5
	SleepTime                  = 2 * time.Second
)
