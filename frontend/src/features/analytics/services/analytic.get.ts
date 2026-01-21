import { apiFetch } from "@/lib/api-fetch";
import { AnalyticMetricsDto } from "../dtos/response/analytic-metrics.dto";
import { BaseResponse } from "@/lib/response/api-response";

export function getAnalyticMetrics(projectId: number, startTime?: string, endTime?: string) {
  return apiFetch<BaseResponse<AnalyticMetricsDto>>(`/projects/${projectId}/analytics`, {
    query: {
      startTime: startTime,
      endTime: endTime,
    },
    withCredentials: true,
  });
}