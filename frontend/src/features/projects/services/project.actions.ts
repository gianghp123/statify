'use server'

import { apiFetch } from "@/lib/api-fetch";
import { CreateProjectRequestDto, UpdateProjectRequestDto } from "../dtos/request/project.request.dto";
import { ProjectDto } from "../dtos/response/project.response.dto";
import { BaseResponse } from "@/lib/response/api-response";

export async function createProject(data: CreateProjectRequestDto) {
  return apiFetch<BaseResponse<ProjectDto>>("/projects", {
    method: "POST",
    body: JSON.stringify(data),
    withCredentials: true,
  });
}

export async function updateProject(id: number, data: UpdateProjectRequestDto) {
  return apiFetch<BaseResponse<ProjectDto>>(`/projects/${id}`, {
    method: "PATCH",
    body: JSON.stringify(data),
    withCredentials: true,
  });
}

export async function deleteProject(id: number) {
  return apiFetch<BaseResponse<void>>(`/projects/${id}`, {
    method: "DELETE",
    withCredentials: true,
  });
}
