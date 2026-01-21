import { z } from "zod"

export interface AnalyticMetricsDto {
  performanceMetrics: PerformanceMetricsDTO;
  projectOverviewMetrics: ProjectOverviewMetricsDTO;
  timeSeriesPoints: TimeSeriesPointDTO[];
}

export interface PerformanceMetricsDTO {
  avgResponseMs: number;
}

export interface ProjectOverviewMetricsDTO {
  totalRequests: number;
  totalBandwidth: number;
  avgResponseMs: number;
  errorRatePercent: number;
}

export interface TimeSeriesPointDTO {
  timestamp: Date;
  requests: number;
  bandwidth: number;
}


export const AnalyticMetricsSchema = z.object({
  // 1. Match the SNAKE_CASE keys coming from the API
  performance_metrics: z.object({
    avg_response_ms: z.coerce.number(),
  }),
  project_overview_metrics: z.object({
    total_requests: z.coerce.number(),
    total_bandwidth: z.coerce.number(),
    avg_response_ms: z.coerce.number(),
    error_rate_percent: z.coerce.number(),
  }),
  time_series_points: z.array(
    z.object({
      // Based on your logs, these keys seem to match, but 'timestamp' 
      // often comes in as a string so we keep coerce.date()
      timestamp: z.coerce.date(),
      requests: z.coerce.number(),
      bandwidth: z.coerce.number(),
    })
  ),
}).transform((data) => ({
  performanceMetrics: {
    avgResponseMs: data.performance_metrics.avg_response_ms,
  },
  projectOverviewMetrics: {
    totalRequests: data.project_overview_metrics.total_requests,
    totalBandwidth: data.project_overview_metrics.total_bandwidth,
    avgResponseMs: data.project_overview_metrics.avg_response_ms,
    errorRatePercent: data.project_overview_metrics.error_rate_percent,
  },
  timeSeriesPoints: data.time_series_points,
}));