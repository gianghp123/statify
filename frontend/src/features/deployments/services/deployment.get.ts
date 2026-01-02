import { apiFetch } from "@/lib/api-fetch";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";

export async function getDeployments(projectId: number, page: number = 1, limit: number = 10) {
  return apiFetch<DeploymentDto[]>(`/projects/${projectId}/deployments`, {
    query: { page, limit },
    withCredentials: true,
  });
}

export async function getDeploymentStatus(projectId: number, deploymentId: number) {
  return apiFetch<DeploymentDto>(`/projects/${projectId}/deployments/${deploymentId}`, {
    withCredentials: true,
  });
}
