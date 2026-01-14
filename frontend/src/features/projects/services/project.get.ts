import { apiFetch } from "@/lib/api-fetch";
import { ProjectDto } from "../dtos/response/project.response.dto";
import { BasePaginatedResponse, BaseResponse } from "@/lib/response/api-response";

export async function getProjects(page: number = 1, limit: number = 10, status: string = '') {
  return apiFetch<BasePaginatedResponse<ProjectDto[]>>("/projects", {
    query: { page, limit, status },
    withCredentials: true,
  });
}

export async function getProjectById(projectId: number) {
  return apiFetch<BaseResponse<ProjectDto>>(`/projects/${projectId}`, {
    withCredentials: true,
  });
}
