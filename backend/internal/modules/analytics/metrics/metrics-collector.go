package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/analytics/repository"
)

type aggregatorBucket struct {
	ProjectID    uint
	DeploymentID uint
	MinuteTs     time.Time

	RequestCount  int64
	BytesServed   int64
	TotalDuration int64

	Status2xx int64
	Status3xx int64
	Status4xx int64
	Status5xx int64
}

type MetricCollector struct {
	eventChan     chan MetricEvent
	metricsRepo   repository.DeploymentMetricsMinuteRepositoryInterface
	aggregatorMap map[string]*aggregatorBucket
}

func NewMetricCollector(metricsRepo repository.DeploymentMetricsMinuteRepositoryInterface) *MetricCollector {
	return &MetricCollector{
		metricsRepo:   metricsRepo,
		eventChan:     make(chan MetricEvent, 100),
		aggregatorMap: make(map[string]*aggregatorBucket),
	}
}

func (c *MetricCollector) EmitMetric(metricEvent MetricEvent) {
	select {
	case c.eventChan <- metricEvent:
	default:
	}
}

func (c *MetricCollector) Start(ctx context.Context) {
	go c.runAggregator(ctx)
}

func (c *MetricCollector) runAggregator(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case s := <-c.eventChan:
			c.aggregate(s)
		case <-ticker.C:
			c.flush(ctx)
		case <-ctx.Done():
			c.flush(ctx)
			return
		}
	}
}

func (c *MetricCollector) flush(ctx context.Context) {
	nowMinute := time.Now().Truncate(time.Minute)

	for key, bucket := range c.aggregatorMap {
		if bucket.MinuteTs.Before(nowMinute) {
			model := models.DeploymentMetricsMinute{
				ProjectID:     bucket.ProjectID,
				DeploymentID:  bucket.DeploymentID,
				MinuteTs:      bucket.MinuteTs,
				RequestCount:  bucket.RequestCount,
				BytesServed:   bucket.BytesServed,
				TotalDuration: bucket.TotalDuration,
				Status2xx:     bucket.Status2xx,
				Status3xx:     bucket.Status3xx,
				Status4xx:     bucket.Status4xx,
				Status5xx:     bucket.Status5xx,
			}
			_ = c.metricsRepo.Save(ctx, &model)
			delete(c.aggregatorMap, key)
		}
	}
}

func (c *MetricCollector) aggregate(metricEvent MetricEvent) {
	minute := metricEvent.Timestamp.Truncate(time.Minute)
	key := c.buildKey(metricEvent.DeploymentID, minute)

	bucket, ok := c.aggregatorMap[key]
	if !ok {
		bucket = &aggregatorBucket{
			ProjectID:    metricEvent.ProjectID,
			DeploymentID: metricEvent.DeploymentID,
			MinuteTs:     minute,
		}
		c.aggregatorMap[key] = bucket
	}

	bucket.RequestCount++
	bucket.BytesServed += metricEvent.BytesServed
	bucket.TotalDuration += metricEvent.DurationMs

	switch {
	case metricEvent.StatusCode >= 200 && metricEvent.StatusCode < 300:
		bucket.Status2xx++
	case metricEvent.StatusCode >= 300 && metricEvent.StatusCode < 400:
		bucket.Status3xx++
	case metricEvent.StatusCode >= 400 && metricEvent.StatusCode < 500:
		bucket.Status4xx++
	case metricEvent.StatusCode >= 500:
		bucket.Status5xx++
	}
}

func (c *MetricCollector) buildKey(deploymentID uint, minuteTs time.Time) string {
	return fmt.Sprintf("%d_%d", deploymentID, minuteTs.Unix())
}
