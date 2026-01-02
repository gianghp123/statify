import { apiFetch } from "@/lib/api-fetch";
import { ProjectDto } from "../dtos/response/project.response.dto";

export async function getProjects(page: number = 1, limit: number = 10) {
  return apiFetch<ProjectDto[]>("/projects", {
    query: { page, limit },
    withCredentials: true,
  });
}

export async function getProjectById(projectId: number) {
  return apiFetch<ProjectDto>(`/projects/${projectId}`, {
    withCredentials: true,
  });
}
