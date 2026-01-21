export function streamAnalytics(projectId: number, startTime?: string, endTime?: string) {
  const params = new URLSearchParams();

  if (startTime) params.append("startTime", startTime);
  if (endTime) params.append("endTime", endTime);

  const eventSource = new EventSource(`/api/sse/project/${projectId}/stream-analytics${params.toString() ? `?${params.toString()}` : ""}`);

  return eventSource;
}