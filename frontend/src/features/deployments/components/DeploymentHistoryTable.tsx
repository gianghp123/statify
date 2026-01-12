'use client'
import { MoreHorizontal, Rocket, History, Trash2, AlertCircle } from "lucide-react";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";
import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useState } from "react";
import { deleteDeployment } from "../services/deployment.actions";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

const getStatusBadge = (status: string) => {
  switch (status) {
    case DeploymentStatus.READY:
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary border border-primary/20">
          <span className="size-1.5 rounded-full bg-primary"></span>
          Ready
        </span>
      );
    case DeploymentStatus.QUEUED:
    case DeploymentStatus.UPLOADED:
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-yellow-400/10 px-2.5 py-0.5 text-xs font-medium text-yellow-400 border border-yellow-400/20">
          <span className="size-1.5 rounded-full bg-yellow-400 animate-pulse"></span>
          Building
        </span>
      );
    case DeploymentStatus.FAILED:
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-red-400/10 px-2.5 py-0.5 text-xs font-medium text-red-400 border border-red-400/20">
          <span className="size-1.5 rounded-full bg-red-400"></span>
          Failed
        </span>
      );
    default:
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-muted/10 px-2.5 py-0.5 text-xs font-medium text-muted-foreground border border-muted/20">
          <span className="size-1.5 rounded-full bg-muted-foreground"></span>
          {status}
        </span>
      );
  }
};

interface DeploymentHistoryTableProps {
  deployments: DeploymentDto[];
  projectId: number;
}

export function DeploymentHistoryTable({ deployments, projectId }: DeploymentHistoryTableProps) {
  const router = useRouter();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedDeployment, setSelectedDeployment] = useState<DeploymentDto | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return "--";
    try {
      const date = new Date(dateStr);
      const now = new Date();
      const diffMs = now.getTime() - date.getTime();
      const diffMins = Math.floor(diffMs / 60000);
      const diffHours = Math.floor(diffMs / 3600000);
      const diffDays = Math.floor(diffMs / 86400000);

      if (diffMins < 60) return `${diffMins}m ago`;
      if (diffHours < 24) return `${diffHours}h ago`;
      return `${diffDays}d ago`;
    } catch {
      return dateStr;
    }
  };

  const calculateDuration = (start?: string, end?: string | null) => {
    if (!start || !end) return "--";
    try {
      const s = new Date(start);
      const e = new Date(end);
      const diffMs = e.getTime() - s.getTime();
      const diffSecs = Math.floor(diffMs / 1000);
      if (diffSecs < 60) return `${diffSecs}s`;
      const mins = Math.floor(diffSecs / 60);
      const secs = diffSecs % 60;
      return `${mins}m ${secs}s`;
    } catch {
      return "--";
    }
  };

  const handleDelete = async () => {
    if (!selectedDeployment) return;
    setIsDeleting(true);
    try {
      const res = await deleteDeployment(projectId, selectedDeployment.id);
      if (res.success) {
        toast.success("Deployment deleted successfully");
        router.refresh();
      } else {
        toast.error(res.message || "Failed to delete deployment");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    } finally {
      setIsDeleting(false);
      setIsDeleteDialogOpen(false);
      setSelectedDeployment(null);
    }
  };

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h3 className="text-foreground font-semibold text-lg flex items-center gap-2">
          <History className="text-primary w-5 h-5 drop-shadow-[0_0_8px_var(--neon-brand-glow)]" />
          Deployment History
        </h3>
      </div>
      <div className="w-full overflow-hidden rounded-xl border border-border bg-card shadow-sm transition-all hover:border-border/80">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="border-b border-border bg-muted/30 text-xs uppercase tracking-wider text-muted-foreground font-semibold">
                <th className="px-6 py-4">Deployment ID</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">Context</th>
                <th className="px-6 py-4">Created</th>
                <th className="px-6 py-4">Duration</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {deployments.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-10 text-center text-muted-foreground bg-muted/5">
                    No deployments found for this project.
                  </td>
                </tr>
              ) : (
                deployments.map((deployment) => (
                  <tr key={deployment.id} className="group hover:bg-accent/5 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-foreground text-xs bg-muted/50 border border-border px-1.5 py-0.5 rounded shadow-sm">
                          {deployment.id}
                        </span>
                        <span className="text-muted-foreground truncate max-w-[150px] text-xs font-medium">
                          {deployment.sourceZipObjectKey?.split('/').pop() || "Manual upload"}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {getStatusBadge(deployment.status)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-foreground font-medium">
                      <div className="flex items-center gap-1.5">
                        <Rocket className="w-4 h-4 text-muted-foreground" />
                        Production
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-muted-foreground">
                      {formatDate(deployment.createdAt)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-muted-foreground">
                      {calculateDuration(deployment.createdAt, deployment.finishedAt)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <button className="text-muted-foreground hover:text-foreground transition-all p-2 rounded-lg hover:bg-accent border border-transparent hover:border-border/30 outline-none">
                            <MoreHorizontal className="w-4 h-4" />
                          </button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="bg-popover border-border shadow-xl">
                          <DropdownMenuItem 
                            className="text-destructive focus:text-destructive focus:bg-destructive/10 cursor-pointer flex items-center gap-2 font-medium"
                            onClick={() => {
                              setSelectedDeployment(deployment);
                              setIsDeleteDialogOpen(true);
                            }}
                          >
                            <Trash2 className="w-4 h-4" />
                            Delete Deployment
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
        <div className="border-t border-border bg-muted/20 px-6 py-3 flex items-center justify-between">
          <span className="text-xs text-muted-foreground font-medium">
            Showing {deployments.length} deployments
          </span>
          <div className="flex items-center gap-2">
            <button className="p-1.5 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent border border-transparent hover:border-border/50 transition-all disabled:opacity-30 disabled:hover:bg-transparent disabled:hover:border-transparent" disabled>
              <span className="text-xs font-bold">Prev</span>
            </button>
            <button className="p-1.5 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent border border-transparent hover:border-border/50 transition-all disabled:opacity-30 disabled:hover:bg-transparent disabled:hover:border-transparent" disabled>
              <span className="text-xs font-bold">Next</span>
            </button>
          </div>
        </div>
      </div>

      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent className="bg-card border-border shadow-2xl">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-foreground flex items-center gap-2 text-xl font-bold">
              <AlertCircle className="w-6 h-6 text-destructive" />
              Delete Deployment?
            </AlertDialogTitle>
            <AlertDialogDescription className="text-muted-foreground text-base">
              Are you sure you want to delete deployment <span className="text-foreground font-mono font-bold bg-muted/50 px-1.5 py-0.5 rounded border border-border">#{selectedDeployment?.id}</span>? 
              This will permanently remove the files from our servers.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="gap-2 sm:gap-0">
            <AlertDialogCancel className="bg-background border-border text-foreground hover:bg-accent font-semibold transition-colors">Cancel</AlertDialogCancel>
            <AlertDialogAction 
              onClick={handleDelete}
              className="bg-destructive text-foreground hover:bg-destructive/90 font-bold shadow-lg shadow-destructive/20 border-none transition-all hover:scale-[1.02]"
            >
              {isDeleting ? "Deleting..." : "Confirm Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
