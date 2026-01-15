'use server'

import { apiFetch } from "@/lib/api-fetch";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";
import { BaseResponse } from "@/lib/response/api-response";

export async function createDeployment(projectId: number, file: File) {
  const formData = new FormData();
  formData.append("file", file);

  return apiFetch<BaseResponse<DeploymentDto>>(`/projects/${projectId}/deployments`, {
    method: "POST",
    body: formData,
    withCredentials: true,
  });
}
export async function deleteDeployment(projectId: number, id: number) {
  return apiFetch<BaseResponse<void>>(`/projects/${projectId}/deployments/${id}`, {
    method: "DELETE",
    withCredentials: true,
  });
}

export async function turnDeploymentLive(projectId: number, id: number) {
  return apiFetch<BaseResponse<void>>(`/projects/${projectId}/deployments/${id}/live`, {
    method: "PUT",
    withCredentials: true,
  });
}

export async function turnDeploymentOffline(projectId: number, id: number) {
  return apiFetch<BaseResponse<void>>(`/projects/${projectId}/deployments/${id}/offline`, {
    method: "PUT",
    withCredentials: true,
  });
}

export async function toggleIsSPAMode(id: number) {
  return apiFetch<BaseResponse<void>>(`/deployments/${id}/toggle-spa-mode`, {
    method: "PUT",
    withCredentials: true,
  });
}
