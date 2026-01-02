export interface ProjectDto {
  id: number;
  name: string;
  url: string;
  status: string;
  lastCommit?: string;
  lastCommitHash?: string;
  updatedAt?: string;
  subdomain?: string;
  userId?: number;
  currentDeploymentId?: number | null;
  createdAt?: string;
}
