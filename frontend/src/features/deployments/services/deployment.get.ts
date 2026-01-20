import { apiFetch } from "@/lib/api-fetch";
import { BasePaginatedResponse, BaseResponse } from "@/lib/response/api-response";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";

export async function getDeployments(projectId: number, page: number = 1, limit: number = 10) {
  return apiFetch<BasePaginatedResponse<DeploymentDto[]>>(`/projects/${projectId}/deployments`, {
    query: { page, limit },
    withCredentials: true,
  });
}

export async function getDeploymentStatus(projectId: number, deploymentId: number) {
  return apiFetch<BaseResponse<DeploymentDto>>(`/projects/${projectId}/deployments/${deploymentId}`, {
    withCredentials: true,
  });
}


export async function getGlobalDeploymentHistory(page: number = 1, limit: number = 10) {
  return apiFetch<BasePaginatedResponse<DeploymentDto[]>>(`/deployments`, {
    query: { page, limit },
    withCredentials: true,
  });
}

