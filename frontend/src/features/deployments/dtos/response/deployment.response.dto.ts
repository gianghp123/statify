import { ProjectDto } from "@/features/projects/dtos/response/project.response.dto";
import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";

export interface DeploymentDto {
  id: number;
  project?: ProjectDto;
  status: DeploymentStatus;
  validationError?: string;
  createdAt: string;
  finishedAt?: string;
}
