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
          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-warning/10 text-warning text-xs font-bold border border-warning/20">
            <span className="w-3 h-3 animate-spin border-2 border-warning border-t-transparent rounded-full" />
            BUILDING
          </span>
        );
      case DeploymentStatus.FAILED:
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-error/10 text-error text-xs font-bold border border-error/20">
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
        return "hover:border-warning/30 hover:shadow-warning/5";
      case DeploymentStatus.FAILED:
        return "hover:border-error/30 hover:shadow-error/5";
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
      className={`group bg-card rounded-2xl p-6 border border-border transition-all shadow-xl hover:shadow-2xl dark:shadow-black/40 cursor-pointer ${getBorderColor(project.status)} hover:scale-[1.01] hover:-translate-y-1`}
    >
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div className="flex items-start gap-5">
          <div className={`size-14 rounded-2xl bg-linear-to-br ${getGradient(project.name)} border border-border/20 flex items-center justify-center shrink-0 shadow-lg`}>
            <span className="font-black text-foreground text-2xl drop-shadow-md">{project.name[0]}</span>
          </div>
          <div className="space-y-1.5">
            <div className="flex items-center gap-3">
              <h3 className="text-xl font-bold text-foreground group-hover:text-primary transition-colors tracking-tight">{project.name}</h3>
              {getStatusDisplay(project.status)}
            </div>
            <Link
              href={`https://${project.url}`}
              onClick={(e) => e.stopPropagation()}
              className="font-mono text-muted-foreground hover:text-primary transition-all text-xs flex items-center gap-1.5 bg-muted/30 px-2 py-1 rounded-md border border-border/50 w-fit hover:border-primary/30"
            >
              <ExternalLink className="w-3.5 h-3.5 opacity-70" />
              {project.url}
            </Link>
          </div>
        </div>
        <div className="flex items-center gap-6 text-sm text-muted-foreground ml-[76px] md:ml-0">
          <div className="flex items-center gap-3">
            <GitCommit className="w-4 h-4 text-primary/70" />
            <span className="font-mono font-semibold text-foreground/80">{project.lastCommit || "main"}</span>
            <span className="font-mono text-[10px] font-bold bg-muted border border-border px-2 py-0.5 rounded shadow-sm">
              {project.lastCommitHash?.slice(0, 7) || "-------"}
            </span>
          </div>
          <div className="hidden sm:block font-medium">{formatDate(project.updatedAt)}</div>
          <button className="text-muted-foreground hover:text-foreground p-2 rounded-lg hover:bg-muted border border-transparent hover:border-border/50 transition-all outline-none">
            <MoreVertical className="w-5 h-5" />
          </button>
        </div>
      </div>

      {(project.status === DeploymentStatus.QUEUED || project.status === DeploymentStatus.UPLOADED) && (
        <div className="mt-6 ml-[76px] h-1.5 w-full max-w-[250px] bg-muted rounded-full overflow-hidden border border-border/50">
          <div className="h-full bg-warning w-2/3 rounded-full animate-pulse shadow-[0_0_8px_var(--warning)]"></div>
        </div>
      )}

      {project.status === DeploymentStatus.FAILED && (
        <div className="mt-4 ml-[76px] flex items-center gap-2 text-xs text-destructive font-bold bg-destructive/5 w-fit px-3 py-1.5 rounded-lg border border-destructive/20">
          <AlertCircle className="w-4 h-4" />
          <span>Something went wrong with the deployment</span>
        </div>
      )}
    </div>
  );
}
