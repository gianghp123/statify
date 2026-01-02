import { ProjectDto } from "@/features/projects/dtos/response/project.dto";
import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";

export const mockProjects: (ProjectDto & { status: DeploymentStatus; url: string; lastCommit: string; lastCommitHash: string; updatedAt: string })[] = [
  {
    id: 1,
    name: "Portfolio Site",
    subdomain: "portfolio",
    userId: 1,
    currentDeploymentId: 101,
    createdAt: new Date().toISOString(),
    status: DeploymentStatus.READY,
    url: "portfolio.statify.app",
    lastCommit: "main",
    lastCommitHash: "e4f92a",
    updatedAt: "2 days ago",
  },
  {
    id: 2,
    name: "Client Dashboard",
    subdomain: "client-dash",
    userId: 1,
    currentDeploymentId: 102,
    createdAt: new Date().toISOString(),
    status: DeploymentStatus.QUEUED, // Showing as BUILDING in UI
    url: "client-dash.statify.app",
    lastCommit: "dev",
    lastCommitHash: "88b12c",
    updatedAt: "5 mins ago",
  },
  {
    id: 3,
    name: "Marketing Landing",
    subdomain: "landing-v2",
    userId: 1,
    currentDeploymentId: 103,
    createdAt: new Date().toISOString(),
    status: DeploymentStatus.FAILED,
    url: "landing-v2.statify.app",
    lastCommit: "feature/hero",
    lastCommitHash: "1a2b3c",
    updatedAt: "1 week ago",
  },
];
