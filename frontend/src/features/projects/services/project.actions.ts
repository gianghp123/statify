'use server'

import { apiFetch } from "@/lib/api-fetch";
import { CreateProjectRequestDto, UpdateProjectRequestDto } from "../dtos/request/project.request.dto";
import { ProjectDto } from "../dtos/response/project.response.dto";

export async function createProject(data: CreateProjectRequestDto) {
  return apiFetch<ProjectDto>("/projects", {
    method: "POST",
    body: JSON.stringify(data),
    withCredentials: true,
  });
}

export async function updateProject(id: number, data: UpdateProjectRequestDto) {
  return apiFetch<ProjectDto>(`/projects/${id}`, {
    method: "PATCH",
    body: JSON.stringify(data),
    withCredentials: true,
  });
}

export async function deleteProject(id: number) {
  return apiFetch<void>(`/projects/${id}`, {
    method: "DELETE",
    withCredentials: true,
  });
}
