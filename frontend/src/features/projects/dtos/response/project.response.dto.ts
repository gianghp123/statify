import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";

export interface ProjectDto {
  id: number;
  name: string;
  url: string;
  status: DeploymentStatus;
  updatedAt?: string;
  subdomain?: string;
  createdAt?: string;
}
