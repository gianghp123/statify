"use client";

import { useState } from "react";
import { Globe, Loader2, Save, Trash2, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ProjectDto } from "../dtos/response/project.response.dto";
import { updateProject, deleteProject } from "../services/project.actions";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

interface ProjectSettingsDialogProps {
  project: ProjectDto;
  onSuccess?: () => void;
}

export function ProjectSettingsDialog({ project, onSuccess }: ProjectSettingsDialogProps) {
  const router = useRouter();
  const [name, setName] = useState(project.name);
  const [subdomain, setSubdomain] = useState(project.subdomain || "");
  const [isUpdating, setIsUpdating] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const autoSubdomain = name
    .toLowerCase()
    .replace(/[^a-z0-9]/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "") || "project";

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsUpdating(true);
    try {
      const res = await updateProject(project.id, { name, subdomain });
      if (res.success) {
        toast.success("Project updated successfully");
        onSuccess?.();
        router.refresh();
      } else {
        toast.error(res.message || "Failed to update project");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    } finally {
      setIsUpdating(false);
    }
  };

  const handleDelete = async () => {
    setIsDeleting(true);
    try {
      const res = await deleteProject(project.id);
      if (res.success) {
        toast.success("Project deleted successfully");
        router.push("/dashboard");
      } else {
        toast.error(res.message || "Failed to delete project");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="space-y-8 py-4">
      <form onSubmit={handleUpdate} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="edit-name" className="text-foreground font-semibold">Project Name</Label>
          <Input
            id="edit-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="bg-background border-border text-foreground focus:ring-primary h-10"
            required
            disabled={isUpdating || isDeleting}
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <Label htmlFor="edit-subdomain" className="text-foreground font-semibold">Subdomain</Label>
            <Button
              type="button"
              variant="link"
              className="h-auto p-0 text-xs text-primary"
              onClick={() => setSubdomain(autoSubdomain)}
              disabled={isUpdating || isDeleting}
            >
              Reset to suggested
            </Button>
          </div>
          <Input
            id="edit-subdomain"
            value={subdomain}
            onChange={(e) => setSubdomain(e.target.value)}
            className="bg-background border-border text-foreground focus:ring-primary h-10"
            required
            disabled={isUpdating || isDeleting}
          />
          <p className="text-xs text-muted-foreground flex items-center gap-1.5 mt-1">
            <Globe className="w-3.5 h-3.5" />
            URL: <span className="text-foreground font-mono">{project.url}</span>
          </p>
        </div>

        <Button
          type="submit"
          className="w-full bg-primary text-primary-foreground font-bold shadow-neon hover:shadow-neon-strong transition-all"
          disabled={isUpdating || isDeleting || (name === project.name && subdomain === project.subdomain)}
        >
          {isUpdating ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Saving...
            </>
          ) : (
            <>
              <Save className="w-4 h-4 mr-2" />
              Save Changes
            </>
          )}
        </Button>
      </form>

      <div className="pt-6 border-t border-border">
        <div className="bg-red-500/5 border border-red-500/20 rounded-lg p-4 space-y-4">
          <div className="flex items-start gap-3">
            <AlertTriangle className="w-5 h-5 text-red-500 mt-0.5" />
            <div>
              <h4 className="text-red-500 font-bold text-sm">Danger Zone</h4>
              <p className="text-xs text-muted-foreground mt-0.5">
                Deleting this project will permanently remove all associated deployments and data. This action cannot be undone.
              </p>
            </div>
          </div>

          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button
                variant="destructive"
                className="w-full font-bold border border-red-500/30 hover:bg-red-500 hover:text-foreground transition-all"
                disabled={isUpdating || isDeleting}
              >
                <Trash2 className="w-4 h-4 mr-2" />
                Delete Project
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent className="bg-card border-border">
              <AlertDialogHeader>
                <AlertDialogTitle className="text-foreground">Are you absolutely sure?</AlertDialogTitle>
                <AlertDialogDescription className="text-muted-foreground">
                  This will permanently delete the project <span className="text-foreground font-semibold">"{project.name}"</span> and all its deployments.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel className="bg-background border-border text-foreground hover:bg-accent">Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleDelete}
                  className="bg-red-500 text-foreground hover:bg-red-600 font-bold"
                >
                  {isDeleting ? "Deleting..." : "Confirm Delete"}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </div>
  );
}
