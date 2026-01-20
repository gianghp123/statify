export function streamDeploymentStatus() {
  const eventSource = new EventSource(`/api/sse/deployment-status`);
  return eventSource;
}