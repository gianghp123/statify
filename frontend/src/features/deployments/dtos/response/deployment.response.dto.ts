import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";

export interface DeploymentDto {
  id: number;
  projectId: number;
  status: DeploymentStatus;
  outputPrefix: string | null;
  sourceZipObjectKey: string | null;
  validationError: string | null;
  createdAt: string;
  finishedAt: string | null;
}
