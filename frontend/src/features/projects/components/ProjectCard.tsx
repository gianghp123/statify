'use client'

import { ProjectDto } from "../dtos/response/project.response.dto";
import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";
import { ExternalLink, GitCommit, MoreVertical, AlertCircle } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface ProjectCardProps {
  project: ProjectDto;
}

export function ProjectCard({ project }: ProjectCardProps) {
  const router = useRouter();

  const getStatusDisplay = (status: DeploymentStatus) => {
    switch (status) {
      case DeploymentStatus.READY:
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-bold border border-primary/20">
            <span className="size-1.5 rounded-full bg-primary shadow-[0_0_5px_rgba(204,255,0,0.8)]"></span>
            LIVE
          </span>
        );
      case DeploymentStatus.QUEUED:
      case DeploymentStatus.UPLOADED:
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-yellow-500/10 text-yellow-400 text-xs font-bold border border-yellow-500/20">
            <span className="w-3 h-3 animate-spin border-2 border-yellow-400 border-t-transparent rounded-full" />
            BUILDING
          </span>
        );
      case DeploymentStatus.FAILED:
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-red-500/10 text-red-400 text-xs font-bold border border-red-500/20">
            <AlertCircle className="w-3 h-3" />
            FAILED
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-muted/10 text-muted-foreground text-xs font-bold border border-muted/20">
            <span className="size-1.5 rounded-full bg-muted-foreground"></span>
            {status.toUpperCase()}
          </span>
        );
    }
  };

  const getBorderColor = (status: string) => {
    switch (status) {
      case DeploymentStatus.READY:
        return "hover:border-primary/50 hover:shadow-[0_0_20px_rgba(204,255,0,0.05)]";
      case DeploymentStatus.QUEUED:
      case DeploymentStatus.UPLOADED:
        return "hover:border-yellow-400/30 hover:shadow-yellow-400/5";
      case DeploymentStatus.FAILED:
        return "hover:border-red-500/30 hover:shadow-red-500/5";
      default:
        return "hover:border-white/10";
    }
  };

  const getGradient = (name: string) => {
    const gradients = [
      "from-indigo-900 to-purple-900",
      "from-indigo-900 to-blue-900",
      "from-purple-900 to-pink-900",
    ];
    const index = name.length % gradients.length;
    return gradients[index];
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return "Unknown date";
    try {
      const date = new Date(dateStr);
      return date.toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric"
      });
    } catch {
      return dateStr;
    }
  };

  return (
    <div 
      onClick={() => router.push(`/projects/${project.id}`)}
      className={`group bg-card rounded-xl p-5 border border-border transition-all shadow-lg shadow-black/20 cursor-pointer ${getBorderColor(project.status)}`}
    >
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-start gap-4">
          <div className={`size-12 rounded-lg bg-gradient-to-br ${getGradient(project.name)} border border-white/10 flex items-center justify-center shrink-0`}>
            <span className="font-bold text-white text-lg">{project.name[0]}</span>
          </div>
          <div>
            <div className="flex items-center gap-3 mb-1">
              <h3 className="text-lg font-semibold text-white group-hover:text-primary transition-colors">{project.name}</h3>
              {getStatusDisplay(project.status)}
            </div>
            <Link
              href={`https://${project.url}`}
              onClick={(e) => e.stopPropagation()}
              className="font-mono text-lavender-bright hover:text-primary transition-colors text-sm flex items-center gap-1 hover:underline decoration-primary/50 underline-offset-4"
            >
              {project.url}
              <ExternalLink className="w-4 h-4 opacity-50" />
            </Link>
          </div>
        </div>
        <div className="flex items-center gap-6 text-sm text-muted-foreground ml-16 md:ml-0">
          <div className="flex items-center gap-2">
            <GitCommit className="w-4 h-4" />
            <span className="font-mono">{project.lastCommit || "main"}</span>
            <span className="font-mono text-xs bg-white/5 px-1.5 py-0.5 rounded">
              {project.lastCommitHash?.slice(0, 7) || "-------"}
            </span>
          </div>
          <div className="hidden sm:block">{formatDate(project.updatedAt)}</div>
          <button className="text-muted-foreground hover:text-white p-1 rounded hover:bg-white/10 transition-colors">
            <MoreVertical className="w-4 h-4" />
          </button>
        </div>
      </div>

      {(project.status === DeploymentStatus.QUEUED || project.status === DeploymentStatus.UPLOADED) && (
        <div className="mt-4 ml-16 h-1 w-full max-w-[200px] bg-white/5 rounded-full overflow-hidden">
          <div className="h-full bg-yellow-400 w-2/3 rounded-full animate-pulse"></div>
        </div>
      )}

      {project.status === DeploymentStatus.FAILED && (
        <div className="mt-3 ml-16 flex items-center gap-2 text-xs text-red-400">
          <AlertCircle className="w-3 h-3" />
          <span>Something went wrong with the deployment</span>
        </div>
      )}
    </div>
  );
}
