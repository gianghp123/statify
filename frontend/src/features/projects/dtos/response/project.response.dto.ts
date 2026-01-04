import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";

export interface ProjectDto {
  id: number;
  name: string;
  url: string;
  status: DeploymentStatus;
  lastCommit?: string;
  lastCommitHash?: string;
  updatedAt?: string;
  subdomain?: string;
  userId?: number;
  currentDeploymentId?: number | null;
  createdAt?: string;
}
