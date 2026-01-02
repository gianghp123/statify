'use server'

import { apiFetch } from "@/lib/api-fetch";
import { CreateProjectRequestDto } from "../dtos/request/project.request.dto";
import { ProjectDto } from "../dtos/response/project.response.dto";

export async function createProject(data: CreateProjectRequestDto) {
  return apiFetch<ProjectDto>("/projects", {
    method: "POST",
    body: JSON.stringify(data),
    withCredentials: true,
  });
}
