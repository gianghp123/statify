import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";

export interface DeploymentStatusResponseDto {
  id: number;
  status: DeploymentStatus;
}