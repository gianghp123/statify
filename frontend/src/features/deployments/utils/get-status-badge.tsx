import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";
import { AlertCircle, CheckCircle2, Trash2 } from "lucide-react";

export const getStatusBadge = (status: DeploymentStatus) => {
  switch (status) {
    case DeploymentStatus.LIVE:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-bold border border-primary/20">
          <span className="size-1.5 rounded-full bg-primary shadow-[0_0_5px_rgba(204,255,0,0.8)]"></span>
          LIVE
        </span>
      );
    case DeploymentStatus.READY:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-bold border border-primary/20">
          <CheckCircle2 className="w-3 h-3" />
          READY
        </span>
      );
    case DeploymentStatus.QUEUED:
    case DeploymentStatus.UPLOADED:
    case DeploymentStatus.PROCESSING:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-warning/10 text-warning text-xs font-bold border border-warning/20">
          <span className="w-3 h-3 animate-spin border-2 border-warning border-t-transparent rounded-full" />
          {status === DeploymentStatus.PROCESSING ? 'PROCESSING' : 'BUILDING'}
        </span>
      );
    case DeploymentStatus.FAILED:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-error/10 text-error text-xs font-bold border border-error/20">
          <AlertCircle className="w-3 h-3" />
          FAILED
        </span>
      );
    case DeploymentStatus.DELETED:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-muted/10 text-muted-foreground text-xs font-bold border border-muted/20 opacity-60">
          <Trash2 className="w-3 h-3" />
          DELETED
        </span>
      );
    default:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-muted/10 text-muted-foreground text-xs font-bold border border-muted/20">
          <span className="size-1.5 rounded-full bg-muted-foreground"></span>
          {status}
        </span>
      );
  }
};
