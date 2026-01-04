import { Info, Code, GitBranch, Globe, Layout, Clock } from "lucide-react";
import { ProjectDto } from "../dtos/response/project.response.dto";

interface ProjectOverviewProps {
  project: ProjectDto;
}

export function ProjectOverview({ project }: ProjectOverviewProps) {
  const formatDate = (dateStr?: string) => {
    if (!dateStr) return "Never";
    try {
      const date = new Date(dateStr);
      const now = new Date();
      const diffMs = now.getTime() - date.getTime();
      const diffMins = Math.floor(diffMs / 60000);
      const diffHours = Math.floor(diffMs / 3600000);
      const diffDays = Math.floor(diffMs / 86400000);

      if (diffMins < 1) return "Just now";
      if (diffMins < 60) return `${diffMins}m ago`;
      if (diffHours < 24) return `${diffHours}h ago`;
      return `${diffDays}d ago`;
    } catch {
      return dateStr;
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 mb-12">
      {/* Left: Metadata Cards */}
      <div className="lg:col-span-4 flex flex-col gap-4">
        <h3 className="text-white font-semibold text-lg flex items-center gap-2">
          <Info className="text-primary w-5 h-5" />
          Project Details
        </h3>
        <div className="grid grid-cols-1 gap-3">
          {/* Card 1 */}
          <div className="flex items-center gap-4 p-4 rounded-xl border border-border bg-card hover:border-white/10 transition-colors group">
            <div className="size-10 rounded-lg bg-background border border-border flex items-center justify-center text-white shrink-0 group-hover:border-primary/30 transition-colors">
              <Code className="w-5 h-5" />
            </div>
            <div className="flex flex-col">
              <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Framework
              </span>
              <span className="text-white font-semibold">Static Site</span>
            </div>
          </div>
          {/* Card 2 */}
          <div className="flex items-center gap-4 p-4 rounded-xl border border-border bg-card hover:border-white/10 transition-colors group">
            <div className="size-10 rounded-lg bg-background border border-border flex items-center justify-center text-white shrink-0 group-hover:border-primary/30 transition-colors">
              <GitBranch className="w-5 h-5" />
            </div>
            <div className="flex flex-col">
              <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Deployment Branch
              </span>
              <span className="text-white font-semibold flex items-center gap-2">
                {project.lastCommit || "main"}
                <span className="size-1.5 rounded-full bg-primary animate-pulse"></span>
              </span>
            </div>
          </div>
          {/* Card 3 */}
          <div className="flex items-center gap-4 p-4 rounded-xl border border-border bg-card hover:border-white/10 transition-colors group">
            <div className="size-10 rounded-lg bg-background border border-border flex items-center justify-center text-white shrink-0 group-hover:border-primary/30 transition-colors">
              <Globe className="w-5 h-5" />
            </div>
            <div className="flex flex-col">
              <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Subdomain
              </span>
              <span className="text-white font-semibold truncate max-w-[150px]">{project.subdomain}.statify.io</span>
            </div>
          </div>
        </div>
      </div>

      {/* Right: Preview Card */}
      <div className="lg:col-span-8 flex flex-col gap-4">
        <h3 className="text-white font-semibold text-lg flex items-center gap-2">
          <Layout className="text-primary w-5 h-5" />
          Live Preview
        </h3>
        <div className="group relative w-full h-full min-h-[320px] rounded-xl border border-border bg-card overflow-hidden shadow-2xl shadow-black/40">
          {/* Browser Chrome */}
          <div className="absolute top-0 left-0 right-0 h-10 bg-background/90 backdrop-blur border-b border-border flex items-center px-4 gap-2 z-10">
            <div className="flex gap-1.5">
              <div className="size-2.5 rounded-full bg-red-500/20 border border-red-500/50"></div>
              <div className="size-2.5 rounded-full bg-yellow-500/20 border border-yellow-500/50"></div>
              <div className="size-2.5 rounded-full bg-green-500/20 border border-green-500/50"></div>
            </div>
            <div className="mx-auto bg-card/50 h-6 px-3 rounded text-[10px] text-muted-foreground flex items-center font-mono w-64 justify-center truncate border border-white/5">
              {project.url}
            </div>
          </div>
          {/* Placeholder Image for Preview */}
          <div className="w-full h-full pt-10 bg-gradient-to-br from-indigo-900/20 via-background to-background flex items-center justify-center text-muted-foreground/20">
             <Layout className="w-16 h-16" />
          </div>
          {/* Overlay Info */}
          <div className="absolute bottom-4 right-4 bg-background/90 backdrop-blur border border-border px-3 py-1.5 rounded-lg flex items-center gap-2 shadow-lg">
            <span className="flex size-2 relative">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
              <span className="relative inline-flex rounded-full size-2 bg-primary"></span>
            </span>
            <span className="text-xs font-medium text-white flex items-center gap-1.5">
              <Clock className="w-3 h-3" />
              Last deployed {formatDate(project.updatedAt)}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
