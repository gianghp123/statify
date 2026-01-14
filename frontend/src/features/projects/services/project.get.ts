import { apiFetch } from "@/lib/api-fetch";
import { ProjectDto } from "../dtos/response/project.response.dto";
import { BasePaginatedResponse, BaseResponse } from "@/lib/response/api-response";

export async function getProjects(page: number = 1, limit: number = 10) {
  return apiFetch<BasePaginatedResponse<ProjectDto[]>>("/projects", {
    query: { page, limit },
    withCredentials: true,
  });
}

export async function getProjectById(projectId: number) {
  return apiFetch<BaseResponse<ProjectDto>>(`/projects/${projectId}`, {
    withCredentials: true,
  });
}
