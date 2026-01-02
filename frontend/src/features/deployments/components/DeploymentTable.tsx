import { Terminal, CheckCircle2, RefreshCw, AlertCircle, Rocket, FlaskConical, MoreHorizontal } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

const mockFullDeployments = [
  {
    id: "8a2b3c",
    project: "Marketing Site",
    projectColor: "bg-indigo-500",
    commitMessage: "Update landing page hero copy",
    status: "Success",
    author: "Alex Developer",
    branch: "main",
    time: "2m ago",
  },
  {
    id: "b4c5d6",
    project: "Dashboard App",
    projectColor: "bg-emerald-500",
    commitMessage: "Fix navigation z-index bug",
    status: "Building",
    author: "Sarah Design",
    branch: "feat/nav-fix",
    time: "15m ago",
  },
  {
    id: "c9d8e7",
    project: "API Gateway",
    projectColor: "bg-orange-500",
    commitMessage: "Refactor authentication middleware",
    status: "Failed",
    author: "Alex Developer",
    branch: "auth-v2",
    time: "1h ago",
  },
  {
    id: "f1a2b3",
    project: "Docs Site",
    projectColor: "bg-blue-400",
    commitMessage: "Add dark mode support",
    status: "Success",
    author: "Mike Frontend",
    branch: "feature/dark-mode",
    time: "3h ago",
  },
];

const getStatusIcon = (status: string) => {
  switch (status) {
    case "Success":
      return (
        <div className="flex items-center justify-center size-8 rounded-full bg-primary/10 border border-primary/20 text-primary shadow-neon-primary">
          <CheckCircle2 className="w-[18px] h-[18px]" />
        </div>
      );
    case "Building":
      return (
        <div className="flex items-center justify-center size-8 rounded-full bg-primary/10 border border-primary/20 text-primary animate-pulse">
          <RefreshCw className="w-[18px] h-[18px] animate-spin" />
        </div>
      );
    case "Failed":
      return (
        <div className="flex items-center justify-center size-8 rounded-full bg-rose-500/10 border border-rose-500/20 text-rose-400 shadow-[0_0_15px_-3px_rgba(251,113,133,0.3)]">
          <AlertCircle className="w-[18px] h-[18px]" />
        </div>
      );
    default:
      return null;
  }
};

export function DeploymentTable() {
  return (
    <div className="flex flex-col overflow-hidden rounded-xl border border-border bg-card/50 backdrop-blur-sm shadow-2xl">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="border-b border-border bg-card text-xs uppercase tracking-wider text-muted-foreground">
            <th className="px-6 py-4 font-semibold w-32">Status</th>
            <th className="px-6 py-4 font-semibold">Commit</th>
            <th className="hidden lg:table-cell px-6 py-4 font-semibold">Project</th>
            <th className="hidden md:table-cell px-6 py-4 font-semibold w-24">Author</th>
            <th className="hidden sm:table-cell px-6 py-4 font-semibold w-32">Time</th>
            <th className="px-6 py-4 font-semibold w-24 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {mockFullDeployments.map((deployment) => (
            <tr key={deployment.id} className="group hover:bg-primary/5 transition-colors">
              <td className="px-6 py-4 align-top">
                {getStatusIcon(deployment.status)}
              </td>
              <td className="px-6 py-4 align-top">
                <div className="flex flex-col gap-1">
                  <p className="text-white text-sm font-medium leading-snug group-hover:text-primary transition-colors cursor-pointer">
                    {deployment.commitMessage}
                  </p>
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-xs text-primary/70 bg-primary/5 px-1.5 py-0.5 rounded border border-primary/10">
                      {deployment.id}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      on <span className="text-foreground/80">{deployment.branch}</span>
                    </span>
                  </div>
                </div>
              </td>
              <td className="hidden lg:table-cell px-6 py-4 align-middle">
                <div className="inline-flex items-center gap-1.5 bg-white/5 border border-white/10 px-3 py-1 rounded-full text-xs font-medium text-slate-200 shadow-sm">
                  <span className={`size-2 rounded-full ${deployment.projectColor}`}></span>
                  {deployment.project}
                </div>
              </td>
              <td className="hidden md:table-cell px-6 py-4 align-middle">
                <div className="size-8 rounded-full bg-primary/20 border border-primary/40 flex items-center justify-center text-[10px] text-primary font-bold">
                  {deployment.author.split(' ').map(n => n[0]).join('')}
                </div>
              </td>
              <td className="hidden sm:table-cell px-6 py-4 align-middle">
                <span className="text-muted-foreground text-sm">{deployment.time}</span>
              </td>
              <td className="px-6 py-4 align-middle text-right">
                <Button variant="ghost" className="opacity-0 group-hover:opacity-100 transition-opacity h-8 px-3 text-xs flex items-center gap-1.5 hover:bg-white/10">
                  <Terminal className="w-3.5 h-3.5" />
                  Logs
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className="border-t border-border bg-card px-6 py-4 flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Showing <span className="text-white font-medium">1-4</span> of <span className="text-white font-medium">248</span> deployments
        </p>
        <div className="flex gap-2">
           <Button variant="outline" className="h-8 px-3 text-xs" disabled>Previous</Button>
           <Button variant="outline" className="h-8 px-3 text-xs">Next</Button>
        </div>
      </div>
    </div>
  );
}
