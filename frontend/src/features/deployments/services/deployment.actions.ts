'use server'

import { apiFetch } from "@/lib/api-fetch";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";

export async function createDeployment(projectId: number, file: File) {
  const formData = new FormData();
  formData.append("file", file);

  return apiFetch<DeploymentDto>(`/projects/${projectId}/deployments`, {
    method: "POST",
    body: formData,
    withCredentials: true,
  });
}
export async function deleteDeployment(projectId: number, id: number) {
  return apiFetch<void>(`/projects/${projectId}/deployments/${id}`, {
    method: "DELETE",
    withCredentials: true,
  });
}

export async function turnDeploymentLive(projectId: number, id: number) {
  return apiFetch<void>(`/projects/${projectId}/deployments/${id}/live`, {
    method: "PUT",
    withCredentials: true,
  });
}

export async function turnDeploymentOffline(projectId: number, id: number) {
  return apiFetch<void>(`/projects/${projectId}/deployments/${id}/offline`, {
    method: "PUT",
    withCredentials: true,
  });
}
