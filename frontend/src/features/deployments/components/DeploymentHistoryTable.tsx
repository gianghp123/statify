'use client'
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
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";
import { AlertCircle, History, Info, MoreHorizontal, Play, Rocket, StopCircle, Trash2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";
import { deleteDeployment, toggleIsSPAMode, turnDeploymentLive, turnDeploymentOffline } from "../services/deployment.actions";
import { calculateDuration, formatDate } from "@/lib/utils/time.utils";
import { getStatusBadge } from "../utils/get-status-badge";
import { Pagination } from "@/lib/response/api-response";
import { Switch } from "@/components/ui/switch";

interface DeploymentHistoryTableProps {
  initialDeployments: DeploymentDto[];
  projectId: number;
  pagination: Pagination;
}

export function DeploymentHistoryTable({ initialDeployments, projectId, pagination }: DeploymentHistoryTableProps) {
  const router = useRouter();
  const [deployments, setDeployments] = useState<DeploymentDto[]>(initialDeployments);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedDeployment, setSelectedDeployment] = useState<DeploymentDto | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    setDeployments(initialDeployments);
  }, [initialDeployments]);


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

  const handleTurnLive = async (deploymentId: number) => {
    try {
      const res = await turnDeploymentLive(projectId, deploymentId);
      if (res.success) {
        toast.success("Deployment is now live");
        setDeployments((prev) =>
          prev.map((d) => {
            if (d.id === deploymentId) {
              return { ...d, status: DeploymentStatus.LIVE };
            }
            if (d.status === DeploymentStatus.LIVE) {
              return { ...d, status: DeploymentStatus.READY };
            }
            return d;
          })
        );
      } else {
        toast.error(res.message || "Failed to turn deployment live");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    }
  };

  const handleTurnOffline = async (deploymentId: number) => {
    try {
      const res = await turnDeploymentOffline(projectId, deploymentId);
      if (res.success) {
        toast.success("Deployment is now offline");
        setDeployments((prev) =>
          prev.map((d) => {
            if (d.id === deploymentId) {
              return { ...d, status: DeploymentStatus.READY };
            }
            return d;
          })
        );
      } else {
        toast.error(res.message || "Failed to turn deployment offline");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    }
  };

  const handleSpaToggle = async (deploymentId: number, checked: boolean) => {
    try {
      const res = await toggleIsSPAMode(deploymentId);
      if (res.success) {
        toast.success("Deployment SPA mode is now " + (checked ? "on" : "off"));
        setDeployments((prev) =>
          prev.map((d) => {
            if (d.id === deploymentId) {
              return { ...d, isSPA: checked };
            }
            return d;
          })
        );
      } else {
        toast.error(res.message || "Failed to toggle SPA mode");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
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
                <th className="px-6 py-4">
                  <div className="flex items-center gap-1.5">
                    SPA
                    <Tooltip>
                      <TooltipTrigger>
                        <Info className="h-3.5 w-3.5 text-muted-foreground/70 hover:text-primary transition-colors cursor-help" />
                      </TooltipTrigger>
                      <TooltipContent className="max-w-[250px] font-normal tracking-normal normal-case">
                        <p>If SPA is on, your application will be considered Single Page Application. Choosing this wrong will cause the deployment not working properly.</p>
                      </TooltipContent>
                    </Tooltip>
                  </div>
                </th>
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
                          {"Manual upload"}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {getStatusBadge(deployment.status)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Switch
                        checked={!!deployment.isSPA}
                        onCheckedChange={(checked) => handleSpaToggle(deployment.id, checked)}
                      />
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
                          {deployment.status === DeploymentStatus.READY && (
                            <DropdownMenuItem
                              className="cursor-pointer flex items-center gap-2 font-medium"
                              onClick={() => handleTurnLive(deployment.id)}
                            >
                              <Play className="w-4 h-4 text-green-500" />
                              Turn Live
                            </DropdownMenuItem>
                          )}
                          {deployment.status === DeploymentStatus.LIVE && (
                            <DropdownMenuItem
                              className="cursor-pointer flex items-center gap-2 font-medium"
                              onClick={() => handleTurnOffline(deployment.id)}
                            >
                              <StopCircle className="w-4 h-4 text-orange-500" />
                              Turn Offline
                            </DropdownMenuItem>
                          )}
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
            Showing {pagination.totalCount} deployments
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
